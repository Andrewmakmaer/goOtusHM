package main

import (
	"fmt"
	"os"

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
