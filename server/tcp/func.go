package tcp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func SendMsg(uid, key, flag string, code byte, msg []byte) {
	buff := WrapCodeByte(code, msg)
	fmt.Printf("%s %s %s:%s", uid, now(), flag, buff.String()[2:])
	C, ok := ConnMap[uid]
	if !ok {
		fmt.Println(" Error User do not exist")
		return
	}
	if C.Key == key {
		fmt.Println(" OK")
		C.Wch <- buff.Bytes()
	}else {
		fmt.Println(" Error Wrong Key")
	}
}

// 测试 定时自动发送数据
func WorkTest(C *Con) {
	time.Sleep(3 * time.Second)
	fmt.Println("Push MSG: hello")
	C.Wch <- WrapCodeString(ReqMessage, "hello").Bytes()

	time.Sleep(20 * time.Second)
	fmt.Println("Push MSG: world")
	C.Wch <- WrapCodeString(ReqMessage, "world").Bytes()
}

// 测试 命令行获取数据并发送
func TSend(uid string) {
	for {
		msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("TSend",err.Error())
			return
		}
		msg = strings.TrimSpace(msg)
		if msg == "quit" {
			ConnMap[uid].Dch <- true
			return
		}
		ConnMap[uid].Wch <- WrapCodeString(ReqMessage, msg).Bytes()
	}
}

