package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Username string `yaml:"username"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
		Sslmode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

func ReadConfig() *Config {
	var config *Config
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		logrus.Fatalf("Err read config.yaml file - %s", err)
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		logrus.Fatalf("Err unmarshal config.yaml to struct - %s", err)
	}
	return config
}
