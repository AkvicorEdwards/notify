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
	Worker       map[string]chan bool
	sync.RWMutex
}

func NewConnect(serverID string, con net.Conn) *Connect {
	return &Connect{
		ServerID:     serverID,
		Connection:   con,
		DataReceived: make(chan []byte, 2),
		DataSend:     make(chan []byte, 2),
		Termination:  make(chan bool),
		Heartbeat:    make(chan bool, 1),
		Worker:       make(map[string]chan bool),
		RWMutex:      sync.RWMutex{},
	}
}
