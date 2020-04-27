package tcp

import (
	"net"
	"sync"
)

type MSG struct {
	App   string `json:"app"`
	Title string `json:"title"`
	Msg   string `json:"msg"`
}

type Connect struct {
	ServerID     string
	Connection   net.Conn
	DataReceived chan []byte
	DataSend     chan []byte
	Heartbeat    chan bool
	Termination  chan bool
	Worker       *WorkerStruct
	sync.RWMutex
}

type WorkerStruct struct {
	Server chan bool
	Sender chan bool
	Receiver chan bool
	Heartbeat chan bool
}

func NewWorker() *WorkerStruct {
	return &WorkerStruct{
		Server:            make(chan bool, 1),
		Sender:            make(chan bool, 1),
		Receiver:          make(chan bool, 1),
		Heartbeat:         make(chan bool, 1),
	}
}

func NewConnect(serverID string, con net.Conn) *Connect {
	return &Connect{
		ServerID:     serverID,
		Connection:   con,
		DataReceived: make(chan []byte, 2),
		DataSend:     make(chan []byte, 2),
		Termination:  make(chan bool),
		Heartbeat:    make(chan bool, 1),
		Worker:       NewWorker(),
		RWMutex:      sync.RWMutex{},
	}
}
