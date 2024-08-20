package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Logger  Logging `yaml:"logging"`
	Brocker Brocker `yaml:"brocker"`
}

type Logging struct {
	Level string `yaml:"level"`
	Type  string `yaml:"type"`
}

type Brocker struct {
	Endpoint string `yaml:"endpoint"`
	Queue    string `yaml:"queue"`
}

func NewConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return Config{}
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
		return Config{}
	}

	return config
}
