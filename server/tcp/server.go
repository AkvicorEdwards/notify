package tcp

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	Err                             byte = 0
	Req                             byte = 1
	Res                             byte = 2
	ReqHEARTBEAT                    byte = 3
	ResHEARTBEAT                    byte = 4
	ReqRegister                     byte = 5
	ResRegistrationFailure          byte = 6
	ResRegistrationFailureUserExist byte = 7
	ResRegistrationSuccessful       byte = 8

	Notify							byte = 20

	MaxMessageLength				int  = 4*200
)


type Con struct {
	Conn net.Conn

	// Received Data
	Rch  chan []byte

	// Send Data
	Wch  chan []byte

	// Send a heartbeat but not respond
	Heart bool

	// Received a heartbeat
	RHch chan bool

	// Send a heartbeat
	WHch chan bool

	// Down signal
	Dch  chan bool

	// Client username
	User string

	// For Close goroutine
	Listener bool
	// goroutine already closed?
	Close map[string]chan bool
}

func NewCon(uid string, con net.Conn) *Con {
	return &Con{
		Conn: con,
		Rch:  make(chan []byte),
		Wch:  make(chan []byte),
		Heart: false,
		RHch: make(chan bool),
		WHch: make(chan bool),
		Dch: make(chan bool),
		User: uid,
		Listener: false,
		Close: map[string]chan bool{
			"handler": make(chan bool),
			"receive": make(chan bool),
			"heartbeat": make(chan bool),
			"listener": make(chan bool),
			"send": make(chan bool),
		},
	}
}

var ConnMap map[string]*Con

func ListenTCP(port string) {
	ConnMap = make(map[string]*Con)
	//listen, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(ip), Port: port})
	tcpAddr, _ := net.ResolveTCPAddr("tcp", port)
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("监听端口失败:", err.Error())
		return
	}
	fmt.Println("已初始化连接，等待客户端连接...")

	go Server(listen)
}

func Server(listen *net.TCPListener) {
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常:", err.Error())
			continue
		}
		fmt.Println("客户端连接来自:", conn.RemoteAddr().String())
		go Handler(conn)
	}
}

func Handler(conn net.Conn) {
	var (
		uid  string
		data = make([]byte, 128)
	)
	for {
		n, err := conn.Read(data)
		if err != nil {
			_=conn.Close()
			return
		}
		if data[0] == ReqRegister {
			uid = string(data[2:n])
			if _, ok := ConnMap[uid]; ok {
				_, _ = conn.Write([]byte{ResRegistrationFailureUserExist, '#'})
				_=conn.Close()
				return
			} else {
				ConnMap[uid] = NewCon(uid, conn)
				fmt.Printf("Register Client: %s\n", uid)
				_, _ = conn.Write([]byte{ResRegistrationSuccessful, '#'})
				break
			}
		} else {
			fmt.Println("ERR")
			_, _ = conn.Write([]byte{Err, '#'})
			_=conn.Close()
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
	case <- ConnMap[uid].Close["handler"]:
		ConnMap[uid].Close["handler"] <- true
		return
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

func Heartbeat(C *Con) {
	for {
		select {
		case <-C.Close["heartbeat"]:
			C.Close["heartbeat"] <- true
			return
		case <-C.WHch:
			//fmt.Println(C.User, "Heart Beat Sent")
			ticker := time.NewTicker(5 * time.Second)
			select {
			case <-ticker.C:
				//fmt.Println(C.User, "Down signal sent")
				C.Dch <- true
			case <-C.RHch:
				C.Heart = false
				//fmt.Println(C.User, "Heart Beat OK")
			}
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
			C.WHch <- true
			C.Wch <- []byte{ReqHEARTBEAT, '#'}
			C.Heart = true
			continue
		}
		switch data[0] {
		case ResHEARTBEAT:
			C.RHch <- true
		case Req:
			C.Rch <- data[2:n]
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
			_,_=C.Conn.Write(d)
		}
	}
}

func SendMsg(uid string, msg []byte) {
	var buff bytes.Buffer
	buff.Write([]byte{Notify, '#'})
	buff.Write(msg)
	fmt.Println(uid, "MSG", buff.String())
	C, ok := ConnMap[uid]
	if !ok {
		return
	}
	C.Wch <- buff.Bytes()
}

func Close(C *Con) {
	select {
	case <-C.Dch:
		C.Listener = true
		C.Close["heartbeat"] <- true
		C.Close["receive"] <- true
		C.Close["handler"] <- true
		C.Close["send"] <- true
	}
	closed := 0
	for {
		select {
		case <-C.Close["handler"]:
			closed++
			//fmt.Println(C.User, "Handler Close")
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
			fmt.Println(C.User, "Down, All Closed")
			_ = C.Conn.Close()
			delete(ConnMap, C.User)
			return
		}
	}
}

func Work(C *Con) {
	time.Sleep(3 * time.Second)
	fmt.Println("Push MSG: hello")
	C.Wch <- []byte{Req, '#', 'h', 'e', 'l', 'l', 'o'}


	time.Sleep(20 * time.Second)
	fmt.Println("Push MSG: world")
	C.Wch <- []byte{Req, '#', 'w', 'o', 'r', 'l', 'd'}
}

func TSend(C *Con) {
	for {
		choice, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("Illegal Input")
			continue
		}
		choice = strings.TrimSpace(choice)
		if choice == "quit" {
			C.Dch <- true
			return
		}
		var buff bytes.Buffer
		buff.Write([]byte{Req, '#'})
		buff.Write([]byte(choice))
		C.Wch <- buff.Bytes()
	}
}
