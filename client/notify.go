package main


import (
	"fmt"
	"notify/tcp"
	"os"
	"time"
)

// [ip] [port] [uid] [key]
func main() {
	fmt.Println("IP:", os.Args[1], "PORT:", os.Args[2], "UID:", os.Args[3], time.Now().Format("2006-01-02 15:04:05"))
	go tcp.ListenTCP(fmt.Sprintf("%s:%s", os.Args[1], os.Args[2]), os.Args[3], os.Args[4]) // ip port uid
	select {}
}
