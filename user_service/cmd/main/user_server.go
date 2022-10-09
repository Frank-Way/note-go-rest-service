package main

import (
	"flag"
	"github.com/Frank-Way/note-go-rest-service/user_service/internal/server"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.yaml", "path to user_service's config path")
}

func main() {
	flag.Parse()

	config := server.NewConfig()

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
