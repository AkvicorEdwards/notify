package shell

import (
	"os/exec"
)

func Exec(s string) {
	cmd := exec.Command("/bin/bash", "-c", s)
	_ = cmd.Run()
}