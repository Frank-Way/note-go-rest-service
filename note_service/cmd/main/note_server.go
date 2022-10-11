package main

import (
	"flag"
	"github.com/Frank-Way/note-go-rest-service/note_service/internal/server"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config.yaml", "path to note_service's config path")
}

func main() {
	flag.Parse()

	config := server.NewConfig(configPath)

	s := server.NewServer(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
