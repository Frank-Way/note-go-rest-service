package server

import (
	"github.com/Frank-Way/note-go-rest-service/internal/auth"
	"github.com/Frank-Way/note-go-rest-service/internal/note"
	"github.com/Frank-Way/note-go-rest-service/internal/note/nerror"
	noteStorage "github.com/Frank-Way/note-go-rest-service/internal/note/storage"
	"github.com/Frank-Way/note-go-rest-service/internal/user"
	userStorage "github.com/Frank-Way/note-go-rest-service/internal/user/storage"
	"github.com/Frank-Way/note-go-rest-service/internal/user/uerror"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strconv"
)

type Server struct {
	config   *Config
	logger   *logrus.Logger
	router   *http.ServeMux
	uHandler *user.Handler
	nHandler *note.Handler
}

func NewServer(config *Config) *Server {
	var logger = logrus.New()
	var uStorage user.Storage
	var nStorage note.Storage
	if config.Storage.Type == "in_memory" {
		uStorage = userStorage.NewInMemoryStorage(logger)
		nStorage = noteStorage.NewInMemoryStorage(logger)
	} else if config.Storage.Type == "redis" {
		uDb, err := strconv.Atoi(config.Storage.Configs.Redis.Db.UserDb)
		if err != nil {
			logger.Fatal(err)
		}
		uStorage, err = userStorage.NewRedisStorage(
			config.Storage.Configs.Redis.Url,
			config.Storage.Configs.Redis.Port,
			config.Storage.Configs.Redis.Password,
			uDb,
			logger)
		if err != nil {
			logger.Fatal(err)
		}
		nDb, err := strconv.Atoi(config.Storage.Configs.Redis.Db.NoteDb)
		if err != nil {
			logger.Fatal(err)
		}
		nStorage, err = noteStorage.NewRedisStorage(
			config.Storage.Configs.Redis.Url,
			config.Storage.Configs.Redis.Port,
			config.Storage.Configs.Redis.Password,
			nDb,
			logger)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		logger.Fatal("unknown storage type specified in config")
	}
	var authService = auth.NewAuthService(logger)
	var uService = user.NewService(authService, uStorage, logger)
	var nService = note.NewService(authService, nStorage, logger)
	return &Server{
		config:   config,
		logger:   logger,
		router:   http.NewServeMux(),
		uHandler: user.NewHandler(uService, logger),
		nHandler: note.NewHandler(nService, logger),
	}
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	var addr = net.JoinHostPort(s.config.Listen.BindIP, s.config.Listen.Port)

	s.logger.Info("starting user server on address: " + addr)

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

	uMiddleware := uerror.Middleware(s.uHandler.Handler)
	nMiddleware := nerror.Middleware(s.nHandler.Handler)

	s.router.Handle("/api/v1/users/", uMiddleware)
	s.router.Handle("/api/v1/users", uMiddleware)

	s.router.Handle("/api/v1/notes/", nMiddleware)
	s.router.Handle("/api/v1/notes", nMiddleware)

	// TODO DELETE DEBUG ENDPOINTS
	//s.router.Handle("/debug/allusers", uMiddleware)
	//s.router.Handle("/debug/allnotes", nMiddleware)
}
