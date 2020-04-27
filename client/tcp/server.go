package tcp

import (
	"fmt"
	"log"
	"net"
	"notify/encryption"
	"notify/shell"
	"os"
	"time"
)

var Connections ConnectionsMap
var retry = make(chan bool)

func ListenTCP(address string, uid string, key string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("ListenTCP recover", err)
		}
	}()
	for {
		Server(address, uid, key)
		select {
		case <-retry:
		}
		go shell.Exec(fmt.Sprintf(`notify-send -t 0 "Notify" `+
			`"PID:[%d]\nDisconnected\nRetry after 7 seconds"`, os.Getpid()))
		time.Sleep(7 * time.Second)
	}
}

func Server(address string, uid string, key string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Server recover", err)
		}
	}()
	tcpAddr, _ := net.ResolveTCPAddr("tcp", address)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(now(), "Failed to connect to the server:", err.Error())
		return
	}
	dataT := make([]byte, 128)
	_, _ = conn.Write(WrapCodeDoubleString(ReqRegistration, uid, key))

	n, err := conn.Read(dataT)
	if err != nil {
		log.Println("Connection Error")
		_ = conn.Close()
		return
	}
	data := dataT[:n]
	var cli *Connect
	switch data[0] {
	case ResRegistrationRefuse:
		log.Println(data[2:3], encryption.Base64ByteToString(data[4:]))
		_ = conn.Close()
		return
	case ResRegistrationAllow:
		// TODO Server ID
		cli = NewConnect("serverID", conn)
		Connections.Store("serverID", cli)
	default:
		_ = conn.Close()
		return
	}

	AddWorker(cli, "server")

	go Terminator(cli)
	go Sender(cli)
	go Receiver(cli)
	go Heartbeat(cli)
	go TerminalMessageSender(cli)

	select {
	case <-cli.Worker["server"]:
		cli.Worker["server"] <- true
	}
}

func Heartbeat(cli *Connect) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Heartbeat recover", err)
		}
	}()

	AddWorker(cli, "heartbeat")

	for {
		ticker := time.NewTicker(30 * time.Second)
		select {
		case <-cli.Worker["heartbeat"]:
			//log.Println("Heartbeat down")
			cli.Worker["heartbeat"] <- true
			return
		case <-ticker.C:
			log.Println("Heartbeat timeout")
			cli.Termination <- true
			continue
		case <-cli.Heartbeat:
			//log.Println("Heartbeat OK")
			cli.DataSend <- WrapCode(ResHeartbeat)
			continue
		}
	}
}

func Sender(cli *Connect) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Sender recover", err)
		}
	}()

	AddWorker(cli, "sender")

	for {
		select {
		case <-cli.Worker["sender"]:
			cli.Worker["sender"] <- true
			return
		case d := <-cli.DataSend:
			err := cli.Connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				log.Printf("SID:[%s] Error sending data: "+
					"unable to set timeout. [%s]\n", cli.ServerID, err.Error())
				continue
			}
			n, err := cli.Connection.Write(d)
			if err != nil {
				log.Printf("SID:[%s] Error sending data: "+
					"failed to send. [%s]\n", cli.ServerID, err.Error())
				continue
			}
			if n != len(d) {
				log.Printf("SID:[%s] Error sending data:"+
					" the length of the sent data does not match. "+
					"%d bytes sent. actual length %d bytes\n",
					cli.ServerID, n, len(d))
				continue
			}
		}
	}
}

func Receiver(cli *Connect) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Receiver recover", err)
		}
	}()

	AddWorker(cli, "receiver")

	dataT := make([]byte, MaxMessageLength)
	down := false
	go func() {
		for {
			select {
			case <-cli.Worker["receiver"]:
				down = true
				return
			}
		}
	}()

	for {
		if down {
			cli.Worker["receiver"] <- true
			return
		}
		err := cli.Connection.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			continue
		}
		n, err := cli.Connection.Read(dataT)
		if err != nil {
			continue
		}
		data := dataT[:n]

		switch data[0] {
		case ReqHeartbeat:
			cli.Heartbeat <- true
		case ReqTerminateTheConnection:
			cli.Termination <- true
		case ReqNotify:
			go HandleNotice(data[2:], cli)
		case ReqTerminalMessage:
			fmt.Print("Terminal:")
			fmt.Println(encryption.Base64ByteToString(data[2:]))
		default:
			fmt.Println("Unknown Message:", data[:])
		}
	}

}

func Terminator(cli *Connect) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Terminator recover", err)
		}
	}()
	select {
	case <-cli.Termination:
		log.Println("kill signal generated")
		for _, v := range cli.Worker {
			v <- true
		}
	}

	//log.Println("kill Signal sent")
	closed := 0
	unclosed := make([]string, 0)
	for k, v := range cli.Worker {
		//log.Println("kill", k)
		ticker := time.NewTicker(10 * time.Second)
		select {
		case <-v:
			closed++
			continue
		case <-ticker.C:
			unclosed = append(unclosed, k)
			continue
		}
	}

	_ = cli.Connection.Close()
	if closed == len(cli.Worker) {
		log.Printf("SID:[%s] Connection closed. "+
			"All threads are terminated\n", cli.ServerID)
	} else {
		log.Printf("SID:[%s] Connection closed. "+
			"The following threads are not terminated:%v\n",
			cli.ServerID, unclosed)
	}

	Connections.Delete(cli.ServerID)
	log.Println("Deleted")
	retry <- true
}

func AddWorker(cli *Connect, worker string) {
	cli.Lock()
	cli.Worker[worker] = make(chan bool, 1)
	cli.Unlock()
}
