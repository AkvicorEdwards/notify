package config

import (
	"fmt"
	"notify/parameter"
)

const Version string = "1.0.0"

const Help string = `
Usage:

	version			print config version

	tcp-addr [.]
	tcp-uuid [.]
	tcp-key [.]
	tcp-cert-file [.]
	tcp-key-file [.]


Default:
  TCP:
	Addr:		:7010
	UUID:       Akvicor
	Key         1234
	Cert:		"./cert/cert.pem"
	Key:		"./cert/key.pem"

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
	AddTCPParser()
}


func AddTCPParser() {
	Args["tcp-addr"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setTCPAddr,
	}
	Args["tcp-uuid"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setTCPUUID,
	}
	Args["tcp-key"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setTCPKey,
	}
	Args["tcp-cert-file"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setTCPCertFile,
	}
	Args["tcp-key-file"] = Arg{
		Size:     1,
		Block:    false,
		Executor: setTCPKeyFile,
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

func setTCPAddr(arg []string) {
	TCP.Addr = arg[1]
}
func setTCPUUID(arg []string) {
	TCP.UUID = arg[1]
}
func setTCPKey(arg []string) {
	TCP.Key = arg[1]
}
func setTCPCertFile(arg []string) {
	TCP.CertFile = arg[1]
}
func setTCPKeyFile(arg []string) {
	TCP.KeyFile = arg[1]
}
