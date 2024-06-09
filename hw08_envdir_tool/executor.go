package main

import (
	"fmt"
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)
	for envName, item := range env {
		if item.NeedRemove {
			os.Unsetenv(envName)
		}
		if len(item.Value) != 0 {
			os.Setenv(envName, item.Value)
		}
	}

	stdout, err := command.Output()

	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	fmt.Print(string(stdout))

	return 0
}
