package repositories

import (
	"context"
	"fmt"
	"github.com/Frank-Way/note-go-rest-service/user_service/internal/uerror"
	"github.com/Frank-Way/note-go-rest-service/user_service/internal/user"
	"github.com/sirupsen/logrus"
	"sync"
)

var _ user.Repository = &inMemoryRepository{}

type inMemoryRepository struct {
	sync.Mutex
	logger *logrus.Logger

	users map[string]user.User
}

func NewInMemoryRepository(logger *logrus.Logger) user.Repository {
	imr := &inMemoryRepository{}
	imr.users = make(map[string]user.User)
	imr.logger = logger
	return imr
}

func (imr *inMemoryRepository) Save(ctx context.Context, user user.User) (string, error) {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("save user to in_memory_repository")
	imr.logger.Debug("check if user exists")
	if _, ok := imr.users[user.Login]; ok {
		imr.logger.Debug("user already exists in memory")
		err := uerror.ErrorDuplicate
		err.DeveloperMessage = fmt.Sprintf("there are user with specified login '%s'", user.Login)
		return "", err
	}
	imr.logger.Debug("generate password hash")
	if err := user.GeneratePasswordHash(); err != nil {
		imr.logger.Debugf("error during password hashing: %v", err)
		return "", err
	}
	imr.users[user.Login] = user
	imr.logger.Debug("user was saved")
	return user.Login, nil
}

func (imr *inMemoryRepository) GetByLogin(ctx context.Context, login string) (user.User, error) {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("get user from in_memory_repository")
	imr.logger.Debugf("find user by login: %s", login)
	u, ok := imr.users[login]
	if ok {
		imr.logger.Debug("user found")
		return u, nil
	} else {
		imr.logger.Debugf("user was not found, login: %s", login)
		err := uerror.ErrorNotFound
		err.DeveloperMessage = fmt.Sprintf("user with login '%s' not found", login)
		return user.User{}, err
	}
}

func (imr *inMemoryRepository) GetAll(ctx context.Context) (user.Users, error) {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("get users from in_memory_repository")
	var res []user.User
	for _, v := range imr.users {
		res = append(res, v)
	}
	imr.logger.Debug("users found")
	return res, nil
}

func (imr *inMemoryRepository) Update(ctx context.Context, user user.User) error {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("update user in in_memory_repository")
	imr.logger.Debugf("find user by login: %s", user.Login)
	u, ok := imr.users[user.Login]
	if ok {
		imr.logger.Debug("user found")
		imr.logger.Debug("update password")
		u.Password = user.Password
		imr.logger.Debug("hashing password")
		if err := u.GeneratePasswordHash(); err != nil {
			imr.logger.Debugf("error during password hashing: %v", err)
			return err
		}
		imr.users[u.Login] = u
		imr.logger.Debug("user updated")
		return nil
	} else {
		imr.logger.Debugf("user was not found, login: %s", user.Login)
		err := uerror.ErrorNotFound
		err.DeveloperMessage = fmt.Sprintf("user with login '%s' not found", user.Login)
		return err
	}
}

func (imr *inMemoryRepository) Delete(ctx context.Context, login string) error {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("delete user from in_memory_repository")
	imr.logger.Debugf("find user by login: %s", login)
	u, ok := imr.users[login]
	if ok {
		imr.logger.Debug("user found")
		delete(imr.users, u.Login)
		imr.logger.Debug("user deleted")
		return nil
	} else {
		imr.logger.Debugf("user was not found, login: %s", login)
		err := uerror.ErrorNotFound
		err.DeveloperMessage = fmt.Sprintf("user with login '%s' not found", login)
		return err
	}
}
