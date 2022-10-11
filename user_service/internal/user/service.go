package user

import (
	"context"
	"github.com/Frank-Way/note-go-rest-service/user_service/internal/uerror"
	"github.com/sirupsen/logrus"
)

var _ Service = &service{}

// TODO DELETE DEBUG METHODS

type Service interface {
	SignUp(ctx context.Context, dto CreateUserDTO) (string, error)
	SignIn(ctx context.Context, login string, dto AuthUserDTO) error
	ChangePassword(ctx context.Context, login string, dto UpdateUserDTO) error
	DeleteUser(ctx context.Context, login string) error
	//TMPGetAllUsers(ctx context.Context) (Users, error)
}

type service struct {
	repository Repository
	logger     *logrus.Logger
}

func NewService(repository Repository, logger *logrus.Logger) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

//func (s service) TMPGetAllUsers(ctx context.Context) (Users, error) {
//	u, err := s.repository.GetAll(ctx)
//	if err != nil {
//		return Users{}, nil
//	}
//	return u, nil
//}

func (s service) SignUp(ctx context.Context, dto CreateUserDTO) (string, error) {
	s.logger.Info("sign up user")
	s.logger.Debug("check if passwords matching")
	if match := dto.Password == dto.RepeatPassword; !match {
		s.logger.Debug("passwords does not match")
		err := uerror.ErrorPasswordsMismatch
		return "", err
	}
	s.logger.Debug("check if user exists")
	u, err := s.repository.GetByLogin(ctx, dto.Login)
	if err == nil {
		s.logger.Debug("user found")
		err := uerror.ErrorDuplicate
		err.Message = "user already exists"
		return "", err
	}
	s.logger.Debug("create user from dto")
	u = NewUser(dto)
	u.IsActive = true
	s.logger.Debug("pass user to repository to save it")
	uri, err := s.repository.Save(ctx, u)
	if err != nil {
		s.logger.Debugf("error during saving user to repository: %v", err)
		return uri, err
	}
	s.logger.Debug("user saved")
	return "", nil
}

func (s service) SignIn(ctx context.Context, login string, dto AuthUserDTO) error {
	s.logger.Info("sign in user")
	s.logger.Debug("check if user exists")
	u, err := s.repository.GetByLogin(ctx, login)
	if err != nil {
		s.logger.Debug("user not found")
		return err
	}
	s.logger.Debug("check if user is active")
	if !u.IsActive {
		s.logger.Debug("user is not active")
		authErr := uerror.ErrorWrongCredentials
		authErr.Message = "user was deleted"
		return authErr
	}
	s.logger.Debug("check password")
	if err = u.CheckPassword(dto.Password); err != nil {
		s.logger.Debug("wrong password")
		authErr := uerror.ErrorWrongCredentials
		authErr.Err = err
		authErr.Message = "wrong password provided"
		return authErr
	}
	// TODO authorization
	s.logger.Debug("user signed in")
	return nil
}

func (s service) ChangePassword(ctx context.Context, login string, dto UpdateUserDTO) error {
	s.logger.Info("change user's password")
	s.logger.Debug("check if user exists")
	u, err := s.repository.GetByLogin(ctx, login)
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
	nU := UpdateUser(u.Login, dto)
	nU.Id = u.Id
	nU.IsActive = u.IsActive
	s.logger.Debug("pass user to repository to update it")
	if err = s.repository.Update(ctx, nU); err != nil {
		s.logger.Debugf("error during updating user in repository: %v", err)
		return err
	}
	s.logger.Debug("user updated")
	return nil
}

func (s service) DeleteUser(ctx context.Context, login string) error {
	s.logger.Info("delete user")
	s.logger.Debug("check if user exists")
	u, err := s.repository.GetByLogin(ctx, login)
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
	s.logger.Debug("pass user to repository to update it")
	if err = s.repository.Update(ctx, u); err != nil {
		s.logger.Debugf("error duting updating user: %v", err)
		return err
	}
	s.logger.Debug("user deleted")
	return nil
}
