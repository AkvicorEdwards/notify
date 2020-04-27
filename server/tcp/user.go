package tcp

import (
	"fmt"
	"strings"
)

// data format: ReqType#data
// Message: URLBase64 With No Padding
func HandleSendToUser(data []byte, cli *Client) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("HandleSendToUser recover", err)
		}
	}()

	info := strings.SplitN(string(data), "#", 1)

	if len(info) < 2 {
		cli.DataSend <- WrapDoubleCodeString(ResNotifyRefuse, IncorrectNotifyReq,
			ErrorIncorrectNotifyReq.Error())
		return
	}


}

func HandleSendToClient(data []byte, cli *Client) {

}
