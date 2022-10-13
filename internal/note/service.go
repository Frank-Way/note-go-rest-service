package note

import (
	"context"
	"github.com/Frank-Way/note-go-rest-service/internal/auth"
	"github.com/Frank-Way/note-go-rest-service/internal/note/nerror"
	"github.com/sirupsen/logrus"
)

var _ Service = &service{}

type Service interface {
	CreateNote(ctx context.Context, auth string, dto CreateNoteDTO) (string, error)
	UpdateNote(ctx context.Context, auth string, id int, dto UpdateNoteDTO) error
	GetNote(ctx context.Context, auth string, id int) (Note, error)
	GetAllNotes(ctx context.Context, auth string) (Notes, error)
	DeleteNote(ctx context.Context, auth string, id int) error
}

type service struct {
	authMw  *auth.Middleware
	storage Storage
	logger  *logrus.Logger
}

func NewService(authSrv auth.Service, storage Storage, logger *logrus.Logger) Service {
	return &service{
		authMw:  auth.NewMiddleware(authSrv, logger),
		storage: storage,
		logger:  logger,
	}
}

func (s service) CreateNote(ctx context.Context, authStr string, dto CreateNoteDTO) (string, error) {
	s.logger.Info("crete note in service")
	s.logger.Debug("parse authStr")
	authLogin, err := s.authMw.CheckAndParse(ctx, authStr)
	if err != nil {
		s.logger.Debug("error during parsing authStr")
		return "", err
	}
	s.logger.Debug("create note from dto")
	n := NewNote(authLogin, dto)
	s.logger.Debug("pass note to storage to create it")
	uri, err := s.storage.Save(ctx, n)
	if err != nil {
		s.logger.Debugf("error during creating note in storage: %v", err)
		return "", err
	}
	s.logger.Debug("note created in service")
	return uri, nil
}

func (s service) UpdateNote(ctx context.Context, authStr string, id int, dto UpdateNoteDTO) error {
	s.logger.Info("update note in service")
	s.logger.Debug("parse authStr")
	authLogin, err := s.authMw.CheckAndParse(ctx, authStr)
	if err != nil {
		s.logger.Debug("error during parsing authStr")
		return err
	}
	s.logger.Debug("check if note exists")
	n, err := s.storage.GetById(ctx, id)
	if err != nil {
		s.logger.Debug("note not found")
		return err
	}
	s.logger.Debug("check if this is user's note")
	if authLogin != n.Author {
		s.logger.Debug("logins mismatch")
		err := nerror.ErrorNoAuth
		err.DeveloperMessage = "attempt to update another user's note"
		return err
	}
	s.logger.Debug("create note from dto")
	nN := UpdateNote(n.Id, authLogin, dto)
	s.logger.Debug("pass note to storage to save it")
	if err = s.storage.Update(ctx, nN); err != nil {
		s.logger.Debugf("error during updating note in storage: %v", err)
		return err
	}
	s.logger.Debug("note updated in service")
	return nil
}

func (s service) GetNote(ctx context.Context, authStr string, id int) (Note, error) {
	s.logger.Info("get note in service")
	s.logger.Debug("parse authStr")
	authLogin, err := s.authMw.CheckAndParse(ctx, authStr)
	if err != nil {
		s.logger.Debug("error during parsing authStr")
		return Note{}, err
	}
	s.logger.Debug("check if note exists")
	n, err := s.storage.GetById(ctx, id)
	if err != nil {
		s.logger.Debug("note not found")
		return Note{}, err
	}
	s.logger.Debug("check if this is user's note")
	if authLogin != n.Author {
		s.logger.Debug("logins mismatch")
		err := nerror.ErrorNoAuth
		err.DeveloperMessage = "attempt to read another user's note"
		return Note{}, err
	}
	s.logger.Debug("return note in service")
	return n, nil
}

func (s service) GetAllNotes(ctx context.Context, authStr string) (Notes, error) {
	s.logger.Info("get notes in service")
	s.logger.Debug("parse authStr")
	authLogin, err := s.authMw.CheckAndParse(ctx, authStr)
	if err != nil {
		s.logger.Debug("error during parsing authStr")
		return Notes{}, err
	}
	s.logger.Debug("get notes from storage")
	n, err := s.storage.GetAll(ctx, authLogin)
	if err != nil {
		s.logger.Debugf("error during get notes in storage: %v", err)
		return Notes{}, err
	}
	s.logger.Debug("return notes in service")
	return n, nil
}

func (s service) DeleteNote(ctx context.Context, authStr string, id int) error {
	s.logger.Info("get note in service")
	s.logger.Debug("parse authStr")
	authLogin, err := s.authMw.CheckAndParse(ctx, authStr)
	if err != nil {
		s.logger.Debug("error during parsing authStr")
		return err
	}
	s.logger.Debug("check if note exists")
	n, err := s.storage.GetById(ctx, id)
	if err != nil {
		s.logger.Debugf("note not found: %v", err)
		return err
	}
	s.logger.Debug("check if this is user's note")
	if authLogin != n.Author {
		s.logger.Debug("logins mismatch")
		err := nerror.ErrorNoAuth
		err.DeveloperMessage = "attempt to delete another user's note"
		return err
	}
	s.logger.Debug("deleting note from storage")
	if err := s.storage.Delete(ctx, id); err != nil {
		s.logger.Debugf("error during deleting note from storage: %v", err)
		return err
	}
	s.logger.Debug("deleted note in service")
	return nil
}
