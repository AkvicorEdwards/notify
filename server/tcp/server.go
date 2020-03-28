package tcp

import (
	"fmt"
	"net"
	"time"
)

var ConnMap map[string]*Con = make(map[string]*Con)

func ListenTCP(port string) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", port)
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("监听端口失败:", err.Error())
		return
	}
	fmt.Println("已初始化连接，等待客户端连接...")

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常:", err.Error())
			continue
		}
		fmt.Println("客户端连接来自:", conn.RemoteAddr().String())
		go Server(conn)
	}
}

func Server(conn net.Conn) {
	var (
		uid  string
		data = make([]byte, 128)
	)
	for {
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println("获取注册数据异常:", err.Error())
			_ = conn.Close()
			return
		}
		if data[0] == ReqRegister {
			uid = string(data[2:n])
			if _, ok := ConnMap[uid]; ok {
				_, _ = conn.Write(WrapCode(ResRegistrationFailureUserExist))
				_ = conn.Close()
				return
			} else {
				ConnMap[uid] = NewCon(uid, conn)
				fmt.Printf("客户端注册成功: %s %s\n", uid, now())
				_, _ = conn.Write(WrapCode(ResRegistrationSuccessful))
				break
			}
		} else {
			fmt.Println("客户端未发送注册请求")
			_, _ = conn.Write(WrapCode(ResRegistrationFailure))
			_ = conn.Close()
			return
		}
	}

	go Close(ConnMap[uid])
	go Send(ConnMap[uid])
	go Receive(ConnMap[uid])
	go Heartbeat(ConnMap[uid])
	go Listener(ConnMap[uid])
	//go Work(C)
	go TSend(ConnMap[uid])

	select {
	case <- ConnMap[uid].Close["server"]:
		ConnMap[uid].Close["server"] <- true
		return
	}
}

func Listener(C *Con) {
	for {
		if C.Listener {
			C.Close["listener"] <- true
			return
		}
		data := make([]byte, MaxMessageLength)
		err := C.Conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			fmt.Println(err)
			continue
		}
		n, err := C.Conn.Read(data)
		if err != nil {
			if C.Heart {
				continue
			}
			// 发送心跳请求
			C.WHch <- true
			// 防止二次发送心跳
			C.Heart = true
			continue
		}
		switch data[0] {
		case ResHeartbeat:
			C.RHch <- true
		case ReqMessage:
			C.Rch <- data[2:n]
		}
	}
}

func Heartbeat(C *Con) {
	for {
		// 等待心跳请求
		select {
		case <-C.Close["heartbeat"]:
			C.Close["heartbeat"] <- true
			return
		case <-C.WHch:
			// 发送心跳请求，等待响应
			C.Wch <- WrapCode(ReqHeartbeat)
			//fmt.Println(C.User, "Heart Beat Sent")
			ticker := time.NewTicker(5 * time.Second)
			select {
			case <-ticker.C:
				C.Dch <- true
				//fmt.Println(C.User, "Down signal sent")
			case <-C.RHch:
				C.Heart = false
				//fmt.Println(C.User, "Heart Beat OK")
			}
		}
	}
}

func Send(C *Con) {
	for {
		select {
		case <-C.Close["send"]:
			C.Close["send"] <- true
			return
		case d := <-C.Wch:
			//fmt.Println(C.User, "Send:", d[0], string(d[2:]))
			_, _ = C.Conn.Write(d)
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
		}
	}
}

func Close(C *Con) {
	select {
	case <-C.Dch:
		//fmt.Println("Down Signal")
		C.Listener = true
		//fmt.Println("1")
		C.Close["heartbeat"] <- true
		//fmt.Println("2")
		C.Close["receive"] <- true
		//fmt.Println("3")
		C.Close["server"] <- true
		//fmt.Println("4")
		C.Close["send"] <- true
	}
	closed := 0
	for {
		select {
		case <-C.Close["server"]:
			closed++
			//fmt.Println(C.User, "Server Close")
		case <-C.Close["receive"]:
			closed++
			//fmt.Println(C.User, "Receive Close")
		case <-C.Close["heartbeat"]:
			closed++
			//fmt.Println(C.User, "Heartbeat Close")
		case <-C.Close["send"]:
			closed++
			//fmt.Println(C.User, "Send Close")
		case <-C.Close["listener"]:
			closed++
			//fmt.Println(C.User, "Listener Close")
		}
		if closed == 5 {
			fmt.Println(C.User, "Down, All Closed", now())
			_ = C.Conn.Close()
			delete(ConnMap, C.User)
			return
		}
	}
}