package server

type Config struct {
	LogLevel string `yaml:"log_level"`
	Listen   struct {
		Type   string `yaml:"type"`
		BindIP string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	} `yaml:"listen"`
	Repository struct {
		Type    string `yaml:"type"`
		Configs struct {
			InMemory struct {
				Attr1 string `yaml:"attr_1"`
			} `yaml:"in_memory"`
			Redis struct {
				Attr1 string `yaml:"attr_1"`
			} `yaml:"redis"`
			Postgres struct {
				Attr1 string `yaml:"attr_1"`
			} `yaml:"postgres"`
		} `yaml:"configs"`
	} `yaml:"repository"`
}

func NewConfig() *Config {
	var c = &Config{
		LogLevel: "info",
	}
	c.Listen.Type = "port"
	c.Listen.BindIP = "127.0.0.1"
	c.Listen.Port = "8000"
	c.Repository.Type = "in_memory"
	return c
}
