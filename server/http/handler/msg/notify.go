package msg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notify/tcp"
	"strconv"
)

func NotifyAndroid(w http.ResponseWriter, r *http.Request) {
	cid := r.FormValue("cid")
	lev := r.FormValue("lev")

	Msg := tcp.MSG{
		App:   filter(r.FormValue("title")),
		Title: getAppName(filter(r.FormValue("app"))),
		Msg:   filter(r.FormValue("msg")),
	}

	data, err := json.Marshal(Msg)
	if err != nil {
		fmt.Println(err)
	}

	id, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return
	}

	var level byte
	switch lev {
	case "0":
		level = tcp.NotifyLowUrgencyLevel
	case "1":
		level = tcp.NotifyNormalUrgencyLevel
	case "2":
		level = tcp.NotifyCriticalUrgencyLevel
	default:
		level = tcp.NotifyNormalUrgencyLevel
	}

	log.Printf("ClientID:[%s] Level:[%s] MSG:[%s]", cid, lev, string(data))

	tcp.SendNotification(tcp.Int64ToClientId(id), tcp.NotifyAndroidMessage,
		level, data)
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
