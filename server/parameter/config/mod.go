package config

import (
	"fmt"
	"notify/parameter"
)

const Version string = "1.0.0"

const Help string = `
Usage:

	version			print config version

	http-title [.]
	http-addr [.]
	http-protocol [.]
	http-cert [.]
	http-key [.]

	tcp-addr [.]
	tcp-cert [.]
	tcp-key [.]

	mysql-user [.]
	mysql-password [.]
	mysql-host [.]
	mysql-dbname [.]
	mysql-charset [.]
	mysql-parse-time [.]
	mysql-loc [.]

	record-uuid [.]
	record-key [.]
	record-ukey [.]


Default:
  HTTP:
	Title:		Akvicor's Notify
	Addr:		:7020
	Protocol	https
	Cert:		"./cert/cert.pem"
	Key:		"./cert/key.pem"
  TCP:
	Addr:		:7010
	Cert:		"./cert/cert.pem"
	Key:		"./cert/key.pem"
  MySQL:
	User:      "root"
	Password:  "password"
	Host:      "localhost"
	DBName:    "notify"
	Charset:   "utf8mb4"
	ParseTime: "true"
	Loc:       "Local"
  Record:
	record-uuid "Akvicor"
	record-key "123"
	record-mkey "123"
`

var Args = parameter.GetArgs()

type Arg = parameter.Arg

func DeleteParseModule() {
	parameter.DeleteFromBaseArgs("config")
}

func AddParseModule() {
	parameter.AddToBaseArgs("config", Arg{
		Block:    true,
		Executor: config,
	})
	Args["version"] = Arg{
		Size:     0,
		Block:    false,
		Executor: printVersion,
	}
	Args["help"] = Arg{
		Size:     0,
		Block:    false,
		Executor: printHelp,
	}
}

func AddAllParser() {
	AddHTTPParser()
	AddTCPParser()
	AddMySQLParser()
	AddRecordParser()
}

func AddMySQLParser() {
	Args["mysql-user"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setMySqlUser,
	}
	Args["mysql-password"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setMySqlPassword,
	}
	Args["mysql-host"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setMySqlHost,
	}
	Args["mysql-dbname"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setMySqlDBName,
	}
	Args["mysql-charset"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setMySqlCharset,
	}
	Args["mysql-parse-time"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setMySqlParseTime,
	}
	Args["mysql-loc"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setMySqlLoc,
	}
}
func AddHTTPParser() {
	Args["http-title"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setHTTPTitle,
	}
	Args["http-addr"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setHTTPAddr,
	}
	Args["http-protocol"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setHTTPProtocol,
	}
	Args["http-cert"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setHTTPCert,
	}
	Args["http-key"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setHTTPKey,
	}
}

func AddTCPParser() {
	Args["tcp-addr"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setTCPAddr,
	}
	Args["tcp-cert"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setTCPCert,
	}
	Args["tcp-key"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setTCPKey,
	}
}

func AddRecordParser() {
	Args["record-uuid"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setRecordUUID,
	}
	Args["record-key"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setRecordKey,
	}
	Args["record-mkey"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setRecordMapApiKey,
	}
}

func config(args []string) {
	parameter.GenericParseArgs(&Args, args[1:])
}

func printHelp([]string) {
	fmt.Print(Help)
}

func printVersion([]string) {
	fmt.Println(Version)
}

func setHTTPTitle(arg []string) {
	HTTP.Title = arg[1]
}
func setHTTPAddr(arg []string) {
	HTTP.Addr = arg[1]
}
func setHTTPProtocol(arg []string) {
	HTTP.Protocol = arg[1]
}
func setHTTPCert(arg []string) {
	HTTP.Cert = arg[1]
}
func setHTTPKey(arg []string) {
	HTTP.Key = arg[1]
}

func setTCPAddr(arg []string) {
	TCP.Addr = arg[1]
}
func setTCPCert(arg []string) {
	TCP.Cert = arg[1]
}
func setTCPKey(arg []string) {
	TCP.Key = arg[1]
}

func setMySqlUser(arg []string) {
	MySQL.User = arg[1]
}
func setMySqlPassword(arg []string) {
	MySQL.Password = arg[1]
}
func setMySqlHost(arg []string) {
	MySQL.Host = arg[1]
}
func setMySqlDBName(arg []string) {
	MySQL.DBName = arg[1]
}
func setMySqlCharset(arg []string) {
	MySQL.Charset = arg[1]
}
func setMySqlParseTime(arg []string) {
	MySQL.ParseTime = arg[1]
}
func setMySqlLoc(arg []string) {
	MySQL.Loc = arg[1]
}

func setRecordUUID(arg []string) {
	Record.UUID = arg[1]
}
func setRecordKey(arg []string) {
	Record.Key = arg[1]
}
func setRecordMapApiKey(arg []string) {
	Record.MapApiKey = arg[1]
}