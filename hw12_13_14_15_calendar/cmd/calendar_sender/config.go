package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
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
	IntQueue string `yaml:"intqueue"`
}

func NewConfig(path string) Config {
	var config Config
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			return Config{}
		}

		err = yaml.Unmarshal(data, &config)
		if err != nil {
			fmt.Println(err)
			return Config{}
		}
	} else {
		config = getFromEnv(config)
	}

	return config
}

func getFromEnv(cnf Config) Config {
	viper.SetEnvPrefix("SENDER")
	viper.AutomaticEnv()

	cnf.Logger.Level = viper.GetString("Log_Level")
	cnf.Logger.Type = viper.GetString("Log_Type")

	cnf.Brocker.Endpoint = viper.GetString("Brocker_Host")
	cnf.Brocker.Queue = viper.GetString("Brocker_Queue")
	cnf.Brocker.IntQueue = viper.GetString("Int_Queue")

	return cnf
}
