package main

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
			"server": make(chan bool),
			"receive": make(chan bool),
			"heartbeat": make(chan bool),
			"listener": make(chan bool),
			"send": make(chan bool),
		},
	}
}

var ConnMap map[string]*Con

func main() {
	ConnMap = make(map[string]*Con)
	go Server("127.0.0.1:6666", os.Args[1])
	select {

	}
}

func Server(address string, uid string) {
	_, ok := ConnMap[uid]
	if ok {
		fmt.Println("uid exist")
		return
	}
	addr, err := net.ResolveTCPAddr("tcp", address)
	conn, err := net.DialTCP("tcp", nil, addr)
	//	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("连接服务端失败:", err.Error())
		return
	}
	data := make([]byte, 128)

	for {
		var buff bytes.Buffer
		buff.Write([]byte{ReqRegister, '#'})
		buff.Write([]byte(uid))
		_,_=conn.Write(buff.Bytes())
		_,err=conn.Read(data)
		if err != nil {
			fmt.Println("Connection Error")
			_=conn.Close()
			return
		}
		if data[0] == ResRegistrationSuccessful {
			fmt.Println("Registration Successful")
			break
		} else if data[0] == ResRegistrationFailureUserExist {
			fmt.Println("User already Connected")
			_=conn.Close()
			return
		}
	}

	ConnMap[uid] = NewCon(uid, conn)
	fmt.Println("Server Connected")

	go Close(ConnMap[uid])
	go Send(ConnMap[uid])
	go Receive(ConnMap[uid])
	go Listener(ConnMap[uid])
	//go Work(ConnMap[uid])
	go TSend(ConnMap[uid])

	select {
	case <- ConnMap[uid].Close["server"]:
		fmt.Println("Close Server", ConnMap[uid].User)
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
			C.Wch <- []byte{ResHEARTBEAT,'#'}
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
		_, err = C.Conn.Read(data)
		if err != nil {
			C.Listener = true
			C.Dch <- true
			continue
		}
		switch data[0] {
		case ReqHEARTBEAT:
			C.RHch <- true
		case Req:
			C.Rch <- data[2:]
		}
	}
}

func Close(C *Con) {
	select {
	case <-C.Dch:
		C.Close["receive"] <- true
		C.Close["send"] <- true
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
		}
		if closed == 3 {
			_ = C.Conn.Close()
			C.Close["server"] <- true
			return
		}
	}
}

func Work(C *Con) {
	time.Sleep(1 * time.Second)
	fmt.Println("Push MSG: 你好")
	var buff bytes.Buffer
	buff.Write([]byte{Req, '#'})
	buff.Write([]byte("你好"))
	C.Wch <- buff.Bytes()

	buff.Reset()
	time.Sleep(17 * time.Second)
	fmt.Println("Push MSG: world")
	buff.Write([]byte{Req, '#'})
	buff.Write([]byte("world"))
	C.Wch <- buff.Bytes()
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