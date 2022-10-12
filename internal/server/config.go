package server

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

type Config struct {
	LogLevel string `yaml:"log_level"`
	Listen   struct {
		Type   string `yaml:"type"`
		BindIP string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	} `yaml:"listen"`
	Storage struct {
		Type    string `yaml:"type"`
		Configs struct {
			Redis struct {
				Attr1 string `yaml:"attr_1"`
			} `yaml:"redis"`
		} `yaml:"configs"`
	} `yaml:"storage"`
}

var instance *Config
var once sync.Once

func NewConfig(configPath string) *Config {
	once.Do(func() {
		instance = &Config{}
		yamlFile, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatal(err)
		}

		if err := yaml.Unmarshal(yamlFile, &instance); err != nil {
			log.Fatal(err)
		}
	})
	return instance
}
