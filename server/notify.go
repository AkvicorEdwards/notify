package main

import (
	"fmt"
	"notify/db"
	"notify/http"
	"notify/parameter"
	"notify/parameter/config"
	"notify/tcp"
	"os"
	"time"
)


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
	db.SetDEFAULT(config.MySQL)

	go tcp.ListenTCP(config.TCP.Addr)
	go http.ListenHTTP()
	fmt.Println("ListenAndServe:", config.HTTP.Addr, "TCP:", config.TCP.Addr, time.Now().Format("2006-01-02 15:04:05"), "PID:", os.Getpid())
	select {}
}
