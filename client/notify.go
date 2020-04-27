package main

import (
	"fmt"
	"notify/parameter"
	"notify/parameter/config"
	"notify/tcp"
)

// [ip] [port] [uid] [key]
func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Base recover", err)
		}
	}()
	parameter.AddBasicArgs()
	config.AddParseModule()
	config.AddAllParser()
	parameter.ParseArgs()

	tcp.ListenTCP(config.TCP.Addr, config.TCP.UUID, config.TCP.Key)
	//go tcp.ListenTCP(config.TCP.Addr, config.TCP.UUID, config.TCP.Key)
	//select {}
}
