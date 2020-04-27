package tcp

import "C"
import (
	"encoding/json"
	"fmt"
	"log"
	"notify/encryption"
	"notify/shell"
)

func HandleNotice(data []byte, cli *Connect) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("HandleNotice recover", err)
		}
	}()

	msg := encryption.Base64ByteToString(data[4:])
	Notice(msg, cli, data[0], data[2])
}

func Notice(msg string, cli *Connect, typeCode byte, levelCode byte) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Notice recover", err)
		}
	}()
	lev := ""
	switch levelCode {
	case NotifyLowUrgencyLevel:
		lev = "low"
	case NotifyNormalUrgencyLevel:
		lev = "normal"
	case NotifyCriticalUrgencyLevel:
		lev = "critical"
	default:
		lev = "normal"
	}
	log.Println(cli.ServerID, msg)
	switch typeCode {
	case NotifyAndroidMessage:
		data := &MSG{}
		_ = json.Unmarshal([]byte(msg), data)
		go shell.Exec(fmt.Sprintf(`notify-send -t 0 -u %s "%s - %s" "%s\n"`,
			lev, data.App, data.Title, data.Msg))
	case NotifyNormalMessage:
		go shell.Exec(fmt.Sprintf(`notify-send -t 0 -u %s "Terminal" "%s\n"`,
			lev, msg))
	case NotifyToDoNotice:
	default:
		go shell.Exec(fmt.Sprintf(`notify-send -t 0 -u %s "Default" "%s\n"`,
			lev, msg))
	}

}

func TerminalMessageSender(cli *Connect) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("TerminalMessageSender recover", err)
		}
	}()
	for {
		fmt.Println("Target")
		fmt.Println("1. Team")
		fmt.Println("2. Server")
		target := GetTerminalInput()
		if target == "1" {
			fmt.Println("Enter Team ID:")
			id := GetTerminalInput()
			if id == "exit" {
				continue
			}
			fmt.Println("Enter Team Key:")
			key := GetTerminalInput()
			if key == "exit" {
				continue
			}
			fmt.Println("Request:", WrapCodeDoubleString(ReqEnterTeam, id, key))
			cli.DataSend <- WrapCodeDoubleString(ReqEnterTeam, id, key)
		} else if target == "2" {
			fmt.Println("Message:")
			msg := GetTerminalInput()
			if msg == "exit" {
				continue
			}
			cli.DataSend <- WrapCodeString(ReqTerminalMessage, msg)
		}
	}
}
