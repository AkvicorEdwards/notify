package shell

import (
	"fmt"
	"os/exec"
)

func Exec(s string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Exec recover", err)
		}
	}()
	cmd := exec.Command("/bin/bash", "-c", s)
	_ = cmd.Run()
}