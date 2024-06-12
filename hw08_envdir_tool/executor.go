package main

import (
	"errors"
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	for envName, item := range env {
		if item.NeedRemove {
			os.Unsetenv(envName)
		}
		if len(item.Value) != 0 {
			os.Setenv(envName, item.Value)
		}
	}

	err := command.Run()
	if err != nil {
		var exiterror *exec.ExitError
		if ok := errors.As(err, &exiterror); ok {
			return exiterror.ExitCode()
		}

		return 1
	}

	return 0
}
