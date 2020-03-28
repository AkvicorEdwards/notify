package msg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notify/tcp"
	"os"
)

type MSG struct {
	App string `json:"app"`
	Title string `json:"title"`
	Msg string `json:"msg"`
}

func Msg(w http.ResponseWriter, r *http.Request) {
	Msg := MSG{
		App:   filter(r.FormValue("title")),
		Title: getAppName(filter(r.FormValue("app"))),
		Msg:   filter(r.FormValue("msg")),
	}
	data, err := json.Marshal(Msg)
	if err != nil {
		fmt.Println(err)
	}
	tcp.SendMsg(os.Args[3], data)
}

func getAppName(app string) string {
	switch app {
	case "com.tencent.mm":
		return "WeChat"
	case "com.tencent.mobileqq":
		return "QQ"
	case "com.tencent.tim":
		return "TIM"
	case "com.google.android.gm":
		return "Gmail"
	default:
		return app
	}
}

func filter(msg string) string {
	switch msg {
	case "%evtprm1":
		return "APP"
	case "%evtprm2":
		return "TITLE"
	case "%evtprm3":
		return "MSG"
	default:
		return msg
	}
}
