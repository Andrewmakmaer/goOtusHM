package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func getEnvFile(dir string) (string, error) {
	var result string
	file, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		result = scanner.Text()
	} else {
		result = ""
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	result = strings.TrimRight(result, " \t")
	result = strings.ReplaceAll(result, "\x00", "\n")
	return result, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envs := make(Environment)
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range files {
		if strings.ContainsAny(e.Name(), "=") {
			continue
		}
		envVal, err := getEnvFile(dir + "/" + e.Name())
		if err != nil {
			return nil, err
		}

		env := EnvValue{Value: envVal}
		_, ok := os.LookupEnv(e.Name())
		if ok {
			env.NeedRemove = true
		} else {
			env.NeedRemove = false
		}
		envs[e.Name()] = env
	}

	return envs, nil
}
