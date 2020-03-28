package tcp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func Work(C *Con) {
	time.Sleep(1 * time.Second)
	fmt.Println("Push MSG: 你好")
	C.Wch <- WrapCodeString(ReqMessage, "你好").Bytes()

	time.Sleep(17 * time.Second)
	fmt.Println("Push MSG: world")
	C.Wch <- WrapCodeString(ReqMessage, "world").Bytes()
}

// 测试 命令行获取数据并发送
func TSend(C *Con) {
	for {
		msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("TSend",err.Error())
			return
		}
		msg = strings.TrimSpace(msg)
		if msg == "quit" {
			C.Dch <- true
			return
		}
		C.Wch <- WrapCodeString(ReqMessage, msg).Bytes()
	}
}