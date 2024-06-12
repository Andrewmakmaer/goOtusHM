package main

import (
	"fmt"
	"os"
)

func main() {
	envs, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err)
	} else {
		RunCmd(os.Args[2:], envs)
	}
}
