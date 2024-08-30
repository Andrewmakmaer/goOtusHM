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
	Logger     Logging       `yaml:"logging"`
	Storage    StorageConfig `yaml:"storage"`
	Server     HTTP          `yaml:"http"`
	GRPCServer GRPC          `yaml:"grpc"`
}

type Logging struct {
	Level string `yaml:"level"`
	Type  string `yaml:"type"`
}

type StorageConfig struct {
	Type     string          `yaml:"type"`
	InMemory *InMemoryConfig `yaml:"inmemory,omitempty"`
	DB       *DBConfig       `yaml:"db,omitempty"`
}

type InMemoryConfig struct{}

type DBConfig struct {
	Endpoint string `yaml:"endpoint"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
}

type HTTP struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type GRPC struct {
	Port string `yaml:"port"`
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
	viper.SetEnvPrefix("CALENDAR")
	viper.AutomaticEnv()

	cnf.Logger.Level = viper.GetString("Log_Level")
	cnf.Logger.Type = viper.GetString("Log_Type")

	cnf.Storage.Type = viper.GetString("DB_Type")

	var db DBConfig
	db.Endpoint = viper.GetString("DB_Endpoint")
	db.Database = viper.GetString("DB_Database")
	db.User = viper.GetString("DB_User")
	db.Pass = viper.GetString("DB_Pass")
	cnf.Storage.DB = &db

	cnf.Server.Host = viper.GetString("HTTP_Host")
	cnf.Server.Port = viper.GetString("HTTP_Port")

	cnf.GRPCServer.Port = viper.GetString("GRPC_Port")

	return cnf
}
