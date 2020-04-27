package config

import "notify/db"

type TCPStruct struct {
	Addr string
	Cert string
	Key  string
}

type HTTPStruct struct {
	Title    string
	Protocol string
	Addr     string
	Cert     string
	Key      string
}

type RecordStruct struct {
	UUID      string
	Key       string
	MapApiKey string
}

var Record = &RecordStruct{
	UUID:      "Akvicor",
	Key:       "123",
	MapApiKey: "123",
}

var TCP = &TCPStruct{
	Addr: ":7010",
	Cert: "./cert/cert.pem",
	Key:  "./cert/key.pem",
}

var HTTP = &HTTPStruct{
	Title:    "Akvicor's Notify",
	Protocol: "https",
	Addr:     ":7020",
	Cert:     "./cert/cert.pem",
	Key:      "./cert/key.pem",
}

var MySQL = &db.MySQL{
	User:      "root",
	Password:  "password",
	Host:      "localhost",
	DBName:    "notify",
	Charset:   "utf8mb4",
	ParseTime: "true",
	Loc:       "Local",
}
