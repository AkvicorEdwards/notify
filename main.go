package main

import (
	"fmt"
)

func main() {
	fin := make(chan bool)



	select {
	case <-fin:
	}
	fmt.Println("Finished")
}
