package main

import (
	"fmt"
	"net/http"
	"notify/handler"
	"notify/tcp"
	"os"
	"time"
)

// [:Web port] [:tcp port] [uid]

func main() {
	handler.ParsePrefix()
	server := http.Server {
		Addr:           os.Args[1],
		Handler:        &handler.MyHandler{},
		ReadTimeout:    1 * time.Minute,
		MaxHeaderBytes: 8<<20,
	}
	fmt.Println("ListenAndServe:", os.Args[1], "TCP:", os.Args[2])
	tcp.ListenTCP(os.Args[2])
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
