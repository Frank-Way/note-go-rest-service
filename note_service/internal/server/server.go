package server

import (
	"github.com/Frank-Way/note-go-rest-service/note_service/internal/nerror"
	"github.com/Frank-Way/note-go-rest-service/note_service/internal/note"
	"github.com/Frank-Way/note-go-rest-service/note_service/internal/note/repositories"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
)

type Server struct {
	config  *Config
	logger  *logrus.Logger
	router  *http.ServeMux
	handler *note.Handler
}

func NewServer(config *Config) *Server {
	var logger = logrus.New()
	var repository note.Repository
	if config.Repository.Type == "in_memory" {
		repository = repositories.NewInMemoryRepository(logger)
	} else {
		logger.Fatal("unknown repository type specified in config")
	}
	return &Server{
		config:  config,
		logger:  logger,
		router:  http.NewServeMux(),
		handler: note.NewHandler(repository, logger),
	}
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	var addr = net.JoinHostPort(s.config.Listen.BindIP, s.config.Listen.Port)

	s.logger.Info("starting note server on address: " + addr)

	s.configureRouter()

	return http.ListenAndServe(addr, s.router)
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	s.logger.Debug("log level set to " + level.String())

	return nil
}

func (s *Server) configureRouter() {
	s.logger.Debug("configuring router")
	middleware := nerror.Middleware(s.handler.Handler)
	s.router.Handle("/api/v1/notes", middleware)
	s.router.Handle("/api/v1/notes/", middleware)
}
