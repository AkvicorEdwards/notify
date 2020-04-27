package http

import (
	"net/http"
	"notify/http/handler"
	"notify/parameter/config"
	"time"
)

func ListenHTTP()  {
	handler.ParsePrefix()
	server := http.Server {
		Addr:           config.HTTP.Addr,
		Handler:        &handler.MyHandler{},
		ReadTimeout:    1 * time.Minute,
		MaxHeaderBytes: 8<<20,
	}

	if err := server.ListenAndServeTLS(config.HTTP.Cert, config.HTTP.Key); err != nil {
		panic(err)
	}
}