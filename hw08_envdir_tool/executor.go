package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	for envName, item := range env {
		if item.NeedRemove {
			os.Unsetenv(envName)
		}
		if len(item.Value) != 0 {
			os.Setenv(envName, item.Value)
		}
	}

	stdout, err := command.Output()

	fmt.Print(string(stdout))
	if err != nil {
		var exiterror *exec.ExitError
		if ok := errors.As(err, &exiterror); ok {
			return exiterror.ExitCode()
		}

		return 1
	}

	return 0
}
