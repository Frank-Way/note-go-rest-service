package server

import (
	"github.com/Frank-Way/note-go-rest-service/note_service/internal/note"
	"github.com/Frank-Way/note-go-rest-service/note_service/internal/note/repositories"
	"github.com/sirupsen/logrus"
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
	return &Server{
		config:  config,
		logger:  logger,
		router:  http.NewServeMux(),
		handler: note.NewHandler(repositories.NewInMemoryRepository(logger), logger),
	}
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	var addr = s.config.Listen.BindIP + ":" + s.config.Listen.Port

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
	s.router.Handle("/api/v1/notes", s.handler)
	s.router.Handle("/api/v1/notes/", s.handler)
}
