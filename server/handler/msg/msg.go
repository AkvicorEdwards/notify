package msg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notify/tcp"
)

type MSG struct {
	App string `json:"app"`
	Title string `json:"title"`
	Msg string `json:"msg"`
}

func Notify(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	key := r.FormValue("key")
	Msg := MSG{
		App:   filter(r.FormValue("title")),
		Title: getAppName(filter(r.FormValue("app"))),
		Msg:   filter(r.FormValue("msg")),
	}
	data, err := json.Marshal(Msg)
	if err != nil {
		fmt.Println(err)
	}
	tcp.SendMsg(uid, key, "Notify", tcp.ReqNotify, data)
}

func getAppName(app string) string {
	switch app {
	case "com.tencent.mm":
		return "WeChat"
	case "com.tencent.mobileqq":
		return "QQ"
	case "com.tencent.tim":
		return "TIM"
	case "com.microsoft.todos":
		return "TODO"
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
	case "%evtprm4":
		return "EvtPrm4"
	case "%evtprm5":
		return "EvtPrm5"
	case "%evtprm6":
		return "EvtPrm6"
	case "%evtprm7":
		return "EvtPrm7"
	default:
		return msg
	}
}
