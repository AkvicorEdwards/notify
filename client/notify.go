package main


import (
	"fmt"
	"notify/tcp"
	"os"
)

// [ip] [port] [uid]
func main() {
	go tcp.ListenTCP(fmt.Sprintf("%s:%s", os.Args[1], os.Args[2]), os.Args[3]) // ip port uid
	select {}
}