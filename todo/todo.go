package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	pid := os.Getpid()
	fmt.Println(pid)
	fmt.Println(os.Getenv("NOTIFY"))
	prc := exec.Command("ps", "-p", strconv.Itoa(pid), "-v")
	out, err := prc.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}
