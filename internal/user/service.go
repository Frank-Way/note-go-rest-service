package user

import (
	"context"
	"github.com/Frank-Way/note-go-rest-service/internal/auth"
	"github.com/Frank-Way/note-go-rest-service/internal/user/uerror"
	"github.com/sirupsen/logrus"
)

var _ Service = &service{}

type Service interface {
	SignUp(ctx context.Context, dto CreateUserDTO) (string, error)
	SignIn(ctx context.Context, login string, dto AuthUserDTO) (string, error)
	ChangePassword(ctx context.Context, authStr string, dto UpdateUserDTO) error
	DeleteUser(ctx context.Context, authStr string, login string) error
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

func (s service) SignUp(ctx context.Context, dto CreateUserDTO) (string, error) {
	s.logger.Info("sign up user")
	s.logger.Debug("check if passwords matching")
	if match := dto.Password == dto.RepeatPassword; !match {
		s.logger.Debug("passwords does not match")
		err := uerror.ErrorPasswordsMismatch
		return "", err
	}
	s.logger.Debug("check if user exists")
	u, err := s.storage.GetByLogin(ctx, dto.Login)
	if err == nil {
		s.logger.Debug("user found")
		err := uerror.ErrorDuplicate
		err.Message = "user already exists"
		return "", err
	}
	s.logger.Debug("create user from dto")
	u = NewUser(dto)
	u.IsActive = true
	s.logger.Debug("pass user to storage to save it")
	uri, err := s.storage.Save(ctx, u)
	if err != nil {
		s.logger.Debugf("error during saving user to storage: %v", err)
		return "", err
	}
	s.logger.Debug("user saved")
	return uri, nil
}

func (s service) SignIn(ctx context.Context, login string, dto AuthUserDTO) (string, error) {
	s.logger.Info("sign in user")
	s.logger.Debug("check if user exists")
	u, err := s.storage.GetByLogin(ctx, login)
	if err != nil {
		s.logger.Debug("user not found")
		return "", err
	}
	s.logger.Debug("check if user is active")
	if !u.IsActive {
		s.logger.Debug("user is not active")
		authErr := uerror.ErrorWrongCredentials
		authErr.Message = "user was deleted"
		return "", authErr
	}
	s.logger.Debug("check password")
	if err = u.CheckPassword(dto.Password); err != nil {
		s.logger.Debug("wrong password")
		authErr := uerror.ErrorWrongCredentials
		authErr.Err = err
		authErr.Message = "wrong password provided"
		return "", authErr
	}
	s.logger.Debug("generate auth token")
	token, err := s.authMw.GetToken(ctx, login)
	if err != nil {
		s.logger.Debugf("error during getting auth token: %v", err)
		return "", err
	}
	s.logger.Debug("user signed in")
	return token, nil
}

func (s service) ChangePassword(ctx context.Context, authStr string, dto UpdateUserDTO) error {
	s.logger.Info("change user's password")
	s.logger.Debug("parse authStr")
	authLogin, err := s.authMw.CheckAndParse(ctx, authStr)
	if err != nil {
		s.logger.Debug("error during parsing authStr")
		return err
	}
	s.logger.Debug("check if user exists")
	u, err := s.storage.GetByLogin(ctx, authLogin)
	if err != nil {
		s.logger.Debug("user not found")
		return err
	}
	s.logger.Debug("check password")
	if err = u.CheckPassword(dto.OldPassword); err != nil {
		s.logger.Debug("wrong password")
		authErr := uerror.ErrorWrongCredentials
		authErr.Err = err
		authErr.Message = "wrong old password provided"
		return authErr
	}
	s.logger.Debug("check if passwords matching")
	if match := dto.NewPassword == dto.RepeatNewPassword; !match {
		s.logger.Debug("new passwords does not match")
		err := uerror.ErrorPasswordsMismatch
		err.Message = "new passwords does not match"
		return err
	}
	s.logger.Debug("check if user is active")
	if !u.IsActive {
		s.logger.Debug("user is not active")
		err := uerror.ErrorWrongCredentials
		err.Message = "user was deleted"
		return err
	}
	s.logger.Debug("create user from dto")
	nU := UpdateUser(u.Id, u.Login, dto)
	nU.IsActive = u.IsActive
	s.logger.Debug("pass user to storage to update it")
	if err = s.storage.Update(ctx, nU); err != nil {
		s.logger.Debugf("error during updating user in storage: %v", err)
		return err
	}
	s.logger.Debug("user updated")
	return nil
}

func (s service) DeleteUser(ctx context.Context, authStr string, login string) error {
	s.logger.Info("delete user")
	s.logger.Debug("parse authStr")
	authLogin, err := s.authMw.CheckAndParse(ctx, authStr)
	if err != nil {
		s.logger.Debug("error during parsing authStr")
		return err
	}
	s.logger.Debug("check if request is authorized")
	if authLogin != login {
		s.logger.Debug("logins mismatch")
		err := uerror.ErrorNoAuth
		err.DeveloperMessage = "attempt to delete another user"
		return err
	}
	s.logger.Debug("check if user exists")
	u, err := s.storage.GetByLogin(ctx, authLogin)
	if err != nil {
		s.logger.Debug("user not found")
		return err
	}
	s.logger.Debug("check if user is active")
	if !u.IsActive {
		s.logger.Debug("user is not active")
		err := uerror.ErrorWrongCredentials
		err.Message = "user was deleted"
		return err
	}
	s.logger.Debug("set user's status to 'not active'")
	u.IsActive = false
	s.logger.Debug("pass user to storage to update it")
	if err = s.storage.Update(ctx, u); err != nil {
		s.logger.Debugf("error duting updating user: %v", err)
		return err
	}
	s.logger.Debug("user deleted")
	return nil
}
