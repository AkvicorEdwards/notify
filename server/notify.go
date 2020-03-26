package main

import (
	"fmt"
	"net/http"
	"notify/handler"
	"notify/tcp"
	"os"
	"strconv"
	"time"
)

// [Web port] [tcp port] [uid]

func main() {
	handler.ParsePrefix()
	server := http.Server {
		Addr:           os.Args[1],
		Handler:        &handler.MyHandler{},
		ReadTimeout:    1 * time.Minute,
		MaxHeaderBytes: 8<<20,
	}
	p, err := strconv.Atoi(os.Args[2])
	if err != nil {
		// handle error
		fmt.Println(err)
		return
	}
	fmt.Println("ListenAndServe:", os.Args[1], "TCP:", p)
	tcp.ListenTCP("127.0.0.1", p)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
