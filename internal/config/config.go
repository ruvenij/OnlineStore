package config

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

type Config struct {
	Port   int    `json:"port"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

var once sync.Once
var instance *Config
var err error

func GetConfig() *Config {
	once.Do(func() {
		instance, err = loadConfig("config.json")
		if err != nil {
			logrus.Fatal(err)
		}
	})

	return instance
}

func loadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logrus.WithError(err).Error("Failed to open config file")
		return nil, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		logrus.WithError(err).Error("Failed to decode config")
		return nil, err
	}

	return &config, nil
}
