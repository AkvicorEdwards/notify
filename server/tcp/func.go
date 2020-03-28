package tcp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func SendMsg(uid string, msg []byte) {
	buff := WrapCodeByte(ReqNotify, msg)
	fmt.Println(uid, now(), "MSG:", buff.String()[2:])
	C, ok := ConnMap[uid]
	if !ok {
		return
	}
	C.Wch <- buff.Bytes()
}

// 测试 定时自动发送数据
func Work(C *Con) {
	time.Sleep(3 * time.Second)
	fmt.Println("Push MSG: hello")
	C.Wch <- WrapCodeString(ReqMessage, "hello").Bytes()


	time.Sleep(20 * time.Second)
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

