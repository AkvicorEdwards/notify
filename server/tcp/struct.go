package tcp

import "net"

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

	Key string

	Down bool
	// goroutine already closed?
	Close map[string]chan bool
}

func NewCon(uid string, key string, con net.Conn) *Con {
	return &Con{
		Conn: con,
		Rch:  make(chan []byte),
		Wch:  make(chan []byte),
		Heart: false,
		RHch: make(chan bool),
		WHch: make(chan bool),
		Dch: make(chan bool),
		User: uid,
		Key: key,
		Down: false,
		Close: map[string]chan bool{
			"server": make(chan bool),
			"receive": make(chan bool),
			"heartbeat": make(chan bool),
			"listener": make(chan bool),
			"send": make(chan bool),
		},
	}
}
