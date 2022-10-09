package main

import (
	"flag"
	"github.com/Frank-Way/note-go-rest-service/user_service/internal/server"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.yaml", "path to user_service's config path")
}

func main() {
	flag.Parse()

	config := server.NewConfig(configPath)

	s := server.NewServer(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
