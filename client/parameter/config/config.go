package config

type TCPStruct struct {
	Addr     string
	UUID     string
	Key      string
	CertFile string
	KeyFile  string
}

var TCP = &TCPStruct{
	Addr:     ":7010",
	UUID:     "Akvicor",
	Key:      "1234",
	CertFile: "./cert/cert.pem",
	KeyFile:  "./cert/key.pem",
}
