package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  Logging       `yaml:"logging"`
	Storage StorageConfig `yaml:"storage"`
	Brocker Brocker       `yaml:"brocker"`
}

type Logging struct {
	Level string `yaml:"level"`
	Type  string `yaml:"type"`
}

type StorageConfig struct {
	Endpoint string `yaml:"endpoint"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
}

type Brocker struct {
	Endpoint string `yaml:"endpoint"`
	Queue    string `yaml:"queue"`
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
	viper.SetEnvPrefix("SCHEDULER")
	viper.AutomaticEnv()

	cnf.Logger.Level = viper.GetString("Log_Level")
	cnf.Logger.Type = viper.GetString("Log_Type")

	cnf.Storage.Endpoint = viper.GetString("DB_Endpoint")
	cnf.Storage.Database = viper.GetString("DB_Database")
	cnf.Storage.User = viper.GetString("DB_User")
	cnf.Storage.Pass = viper.GetString("DB_Pass")

	cnf.Brocker.Endpoint = viper.GetString("Brocker_Host")
	cnf.Brocker.Queue = viper.GetString("Brocker_Queue")

	return cnf
}
