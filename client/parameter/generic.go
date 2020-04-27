package parameter

import (
	"github.com/fatih/color"
)

func GenericParseArgs(GArgs *Args, args []string) {
	if len(args) == 0 {
		return
	}
	if h, ok := (*GArgs)[args[0]]; ok {
		if h.Block {
			h.Executor(args[:])
			return
		}
		if len(args) <= h.Size {
			color.Red("The command [%s] requires more parameters to execute", args[0])
			return
		}
		h.Executor(args[:1+h.Size])
		GenericParseArgs(GArgs, args[1+h.Size:])
	} else {
		color.Red("Skipped: Unknown command [%s]\n", args[0])
		GenericParseArgs(GArgs, args[1:])
	}
}