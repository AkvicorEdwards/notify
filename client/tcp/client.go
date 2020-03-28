package tcp

import (
	"encoding/json"
	"fmt"
	"net"
	"notify/shell"
	"time"
)

var ConnMap map[string]*Con = make(map[string]*Con)

// [ip] [port] [uid]
func ListenTCP(address string, uid string) {
	for {
		Server(address, uid)
		go shell.Exec(fmt.Sprint(`notify-send -t 0 "Notify" "Disconnected\nRetry after 7 seconds"`))
		time.Sleep(7*time.Second)
	}
}

func Server(address string, uid string) {
	_, ok := ConnMap[uid]
	if ok {
		fmt.Println("UID exist")
		return
	}
	tcpAddr, _ := net.ResolveTCPAddr("tcp", address)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println(now(), "连接服务端失败:", err.Error())
		return
	}
	data := make([]byte, 128)

	for {
		_,_=conn.Write(WrapCodeString(ReqRegister, uid).Bytes())
		_,err=conn.Read(data)
		if err != nil {
			fmt.Println("Connection Error")
			_ = conn.Close()
			return
		}
		if data[0] == ResRegistrationSuccessful {
			fmt.Println("Registration Successful")
			break
		} else if data[0] == ResRegistrationFailureUserExist {
			fmt.Println("User already Connected")
			_ = conn.Close()
			return
		}
	}

	ConnMap[uid] = NewCon(uid, conn)
	fmt.Println("Server Connected", now())
	go shell.Exec(fmt.Sprint(`notify-send -t 0 "Notify" "Connected"`))

	go Close(ConnMap[uid])
	go Send(ConnMap[uid])
	go Receive(ConnMap[uid])
	go Notice(ConnMap[uid])
	go Listener(ConnMap[uid])
	//go Work(ConnMap[uid])
	go TSend(ConnMap[uid])

	select {
	case <- ConnMap[uid].Close["server"]:
		fmt.Println("Close Server", ConnMap[uid].User, now())
		delete(ConnMap, ConnMap[uid].User)
		return
	}
}

func Send(C *Con) {
	for {
		select {
		case <-C.Close["send"]:
			C.Close["send"] <- true
			return
		case d := <-C.Wch:
			//fmt.Println("Send:", d[0], string(d[2:]))
			_,_=C.Conn.Write(d)
		}
	}
}

func Notice(C *Con) {
	for {
		select {
		case <-C.Close["notice"]:
			C.Close["notice"] <- true
			return
		case d := <-C.Nty:
			data := &MSG{}
			_ = json.Unmarshal(d, data)
			fmt.Println(C.User, now(), string(d))
			go shell.Exec(fmt.Sprintf(`notify-send -t 0 "%s - %s" "%s\n"`, data.App, data.Title, data.Msg))
		}
	}

}

func Receive(C *Con) {
	for {
		select {
		case <-C.Close["receive"]:
			C.Close["receive"] <- true
			return
		case d := <-C.Rch:
			fmt.Println(C.User, "Receive:", string(d))
		case <-C.RHch:
			//fmt.Println("Heart Beat Received")
			C.Wch <- []byte{ResHeartbeat,'#'}
		}
	}
}

func Listener(C *Con) {
	for {
		if C.Listener {
			C.Close["listener"] <- true
			return
		}
		data := make([]byte, MaxMessageLength)
		err := C.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		if err != nil {
			fmt.Println(err)
			continue
		}
		n, err := C.Conn.Read(data)
		if err != nil {
			C.Listener = true
			C.Dch <- true
			continue
		}
		switch data[0] {
		case ReqNotify:
			C.Nty <- data[2:n]
		case ReqHeartbeat:
			C.RHch <- true
		case ReqMessage:
			C.Rch <- data[2:n]
		}
	}
}

func Close(C *Con) {
	select {
	case <-C.Dch:
		C.Close["receive"] <- true
		C.Close["send"] <- true
		C.Close["notice"] <- true
	}
	closed := 0
	for {
		select {
		case <-C.Close["receive"]:
			closed++
			//fmt.Println(C.User, "Receive Close")
		case <-C.Close["send"]:
			closed++
			//fmt.Println(C.User, "Send Close")
		case <-C.Close["listener"]:
			closed++
			//fmt.Println(C.User, "Listener Close")
		case <-C.Close["notice"]:
			closed++
			//fmt.Println(C.User, "Notice Close")
		}
		if closed == 4 {
			_ = C.Conn.Close()
			C.Close["server"] <- true
			return
		}
	}
}

