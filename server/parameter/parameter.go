package parameter

import (
	"fmt"
	"os"
)

const (
	Version string = "1.0.0"

	Help string = `
Usage:

	config				Config
	version				print version
	name <your name>	print your name
`
)

var BaseArgs = GetArgs()

func AddBasicArgs() {
	BaseArgs["version"] = Arg{
		Size:     0,
		Block:    false,
		Executor: printVersion,
	}
	BaseArgs["help"] = Arg{
		Size:     0,
		Block:    false,
		Executor: printHelp,
	}
	BaseArgs["name"] = Arg{
		Size:     1,
		Block:    false,
		Executor: name,
	}
}

func AddToBaseArgs(parameter string, arg Arg) {
	if _, ok := BaseArgs[parameter]; ok {
		panic(fmt.Sprintf("parameter [%s] already exist", parameter))
		return
	}
	BaseArgs[parameter] = arg
	return
}

func DeleteFromBaseArgs(parameter string) {
	delete(BaseArgs, parameter)
}

func ParseArgs() {
	GenericParseArgs(&BaseArgs, os.Args[1:])
}

func ParseArgsByString(args []string) {
	GenericParseArgs(&BaseArgs, args[1:])
}

func printHelp([]string) {
	fmt.Print(Help)
}

func printVersion([]string) {
	fmt.Println(Version)
}
