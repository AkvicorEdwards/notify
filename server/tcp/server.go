package tcp

import (
	"fmt"
	"log"
	"net"
	"notify/encryption"
	"strings"
	"time"
)

var Connections ConnectionsMap
var Conn ConnMap
var Teams TeamMap

func ListenTCP(address string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("ListenTCP recover", err)
		}
	}()
	addr, _ := net.ResolveTCPAddr("tcp", address)
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println("Failed to listen on port:", err.Error())
		return
	}

	log.Printf("TCP Server is listening [%s], " +
		"waiting for client to connect...\n", address)
	go CliMessageSender()
	TeamLoader()
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("Accept client connection exception:", err.Error())
			continue
		}
		log.Println("Client connection comes from:", conn.RemoteAddr().String())
		go Server(conn)
	}
}

func Server(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Server recover", err)
		}
	}()

	var (
		cli *Client
		data = make([]byte, 128)
		info []string
		ok   bool
		n    int
		err  error
	)

	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		_, _ = conn.Write(WrapDoubleCodeString(ResRegistrationRefuse, Unknown,
			ErrorUnknown.Error()))
		_ = conn.Close()
		return
	}
	n, err = conn.Read(data)
	if err != nil {
		_, _ = conn.Write(WrapDoubleCodeString(ResRegistrationRefuse, Unknown,
			ErrorUnknown.Error()))
		_ = conn.Close()
		return
	}

	if data[0] == ReqRegistration && n > 2 {
		info = strings.Split(string(data[2:n]), "#")
		if len(info) < 2 {
			_, _ = conn.Write(WrapDoubleCodeString(ResRegistrationRefuse,
				IllegalRegistrationInfo, ErrorIllegalRegistrationInfo.Error()))
			_ = conn.Close()
			return
		}
		// uuid
		info[0] = encryption.Base64StringToString(info[0])
		// key
		info[1] = encryption.Base64StringToString(info[1])
		// ip
		ip := strings.Split(conn.RemoteAddr().String(), ":")[0]

		cli, ok = AddConnections(info[0], ip, info[1], conn)
		if ok {
			log.Printf("Client [%s] registered successfully: cid:[%v] ip:[%s]\n",
				info[0], cli.Id, conn.RemoteAddr().String())
		} else {
			_ = conn.Close()
			return
		}
	} else {
		_, _ = conn.Write(WrapDoubleCodeString(ResRegistrationRefuse,
			IllegalRegistrationInfo, ErrorIllegalRegistrationInfo.Error()))
		_ = conn.Close()
		return
	}

	AddWorker(cli, "server")

	go Terminator(cli)
	go Sender(cli)
	go Receiver(cli)
	go Heartbeat(cli)
	go ApiDeadlineTicker(cli)

	select {
	case <-cli.Worker["server"]:
		cli.Worker["server"] <- true
		return
	}
}

func Sender(cli *Client) {
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
				log.Printf("UUID:[%v] CID:[%v] Error sending data:" +
					" unable to set timeout. [%s]\n", cli.UUID, cli.Id, err.Error())
				continue
			}
			n, err := cli.Connection.Write(d)
			if err != nil {
				log.Printf("UUID:[%v] CID:[%v] Error sending data:" +
					" failed to send. [%s]\n", cli.UUID, cli.Id, err.Error())
				continue
			}
			if n != len(d) {
				log.Printf("UUID:[%v] CID:[%v] Error sending data:" +
					" the length of the sent data does not match. " +
					"%d bytes sent. actual length %d bytes\n",
					cli.UUID, cli.Id, n, len(d))
				continue
			}
		}
	}
}

func Receiver(cli *Client) {
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
		err := cli.Connection.SetReadDeadline(time.Now().Add(30 * time.Second))
		if err != nil {
			continue
		}
		n, err := cli.Connection.Read(dataT)
		if err != nil || n < 2 {
			continue
		}

		data := dataT[:n]
		switch data[0] {
		case ResHeartbeat:
			cli.HeartbeatDataReceived <- true
		case ReqTerminateTheConnection:
			cli.Termination <- true
		case ReqEnterTeam:
			go HandleEnterTeam(data[2:], cli)
		case ReqTeamMessage:
			go HandleSendToTeam(data[2:], cli)
		case ReqUserMessage:
			go HandleSendToUser(data[2:], cli)
		case ReqClientMessage:
			go HandleSendToClient(data[2:], cli)
		case ReqTerminalMessage:
			fmt.Printf("UUID:[%v] CID:[%v] Terminal Message:[%s]\n",
				cli.UUID, cli.Id, encryption.Base64ByteToString(data[2:]))
		default:
			fmt.Printf("UUID:[%v] CID:[%v] Terminal Message:[%s]\n",
				cli.UUID, cli.Id, data)
		}
	}

}

func Heartbeat(cli *Client) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Heartbeat recover", err)
		}
	}()
	AddWorker(cli, "heartbeat")

	for {
		//log.Printf("UUID:[%v] CID:[%v] Heartbeat 1 prepared",
		//	cli.UUID, cli.Id)
		ticker := time.NewTicker(5 * time.Second)
		select {
		case <-cli.Worker["heartbeat"]:
			cli.Worker["heartbeat"] <- true
			return
		case <-ticker.C:
			//log.Printf("UUID:[%v] CID:[%v] Heartbeat 2 Send",
			//	cli.UUID, cli.Id)
			cli.DataSend <- WrapCode(ReqHeartbeat)
			beatTicker := time.NewTicker(10 * time.Second)
			select {
			case <-cli.Worker["heartbeat"]:
				cli.Worker["heartbeat"] <- true
				return
			case <-beatTicker.C:
				log.Printf("UUID:[%v] CID:[%v] Heartbeat Failure",
					cli.UUID, cli.Id)
				cli.Termination <- true
			case <-cli.HeartbeatDataReceived:
				//log.Printf("UUID:[%v] CID:[%v] Heartbeat 3 ok",
				//	cli.UUID, cli.Id)
				continue
			}
		}
	}
}

func ApiDeadlineTicker(cli *Client) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("ApiDeadlineTicker recover", err)
		}
	}()
	tk := cli.Deadline - time.Now().Unix()
	if tk <= 0 {
		cli.DataSend <- WrapDoubleCodeString(ResError, IllegalAPIkey,
			ErrorIllegalAPIkey.Error())
		cli.Termination <- true
		return
	}
	AddWorker(cli, "apiDeadlineTicker")
	ApiDeadTicker := time.NewTicker(time.Duration(tk)*time.Second)
	select {
	case <-cli.Worker["apiDeadlineTicker"]:
		cli.Worker["apiDeadlineTicker"] <- true
		return
	case <-ApiDeadTicker.C:
		cli.DataSend <- WrapDoubleCodeString(ResError, IllegalAPIkey,
			ErrorIllegalAPIkey.Error())
		cli.Termination <- true
	}
}

func Terminator(cli *Client) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Terminator recover", err)
		}
	}()
	select {
	case <-cli.Termination:
		log.Printf("UUID:[%v] CID:[%v] Kill Signal Generated",
			cli.UUID, cli.Id)
		for _, v := range cli.Worker {
			v <- true
		}
		cli.Lock()
		for _, v := range cli.ConnectedTeam {
			v.Lock()
			delete(v.ConnectedClient, cli.Id)
			v.Unlock()
		}
		cli.Unlock()
	}
	closed := 0
	unclosed := make([]string, 0)
	for k, v := range cli.Worker {
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
	if closed == len(cli.Worker) {
		log.Printf("UUID:[%v] CID:[%v] Connection closed. " +
			"All threads are terminated", cli.UUID, cli.Id)
	} else {
		log.Printf("UUID:[%v] CID:[%v] Connection closed. " +
			"The following threads are not terminated:%v",
			cli.UUID, cli.Id, unclosed)
	}
	Conn.Delete(cli.Id)

	_ = cli.Connection.Close()
	k, ok := Connections.RLoad(cli.UUID)
	if ok {
		k.Lock()
		delete(k.ConnectedClient, cli.Id)
		k.Unlock()
		log.Printf("UUID:[%v] CID:[%v] Deleted\n", cli.UUID, cli.Id)
	} else {
		log.Printf("UUID:[%v] CID:[%v] Delete Failure\n", cli.UUID, cli.Id)
	}
}

func AddWorker(cli *Client, worker string) {
	cli.Lock()
	cli.Worker[worker] = make(chan bool, 1)
	cli.Unlock()
}
