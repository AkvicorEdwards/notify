package handler

import (
	"net/http"
	"notify/handler/msg"
)

type str2func map[string]func(http.ResponseWriter, *http.Request)

var public str2func

func ParsePrefix() {
	public = make(str2func)

	public["/msg"] = msg.Msg
}

type MyHandler struct {}

func (*MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Public
	if h, ok := public[r.URL.Path]; ok {
		h(w, r)
		return
	}
}



