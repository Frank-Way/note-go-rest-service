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

	users  map[uint]user.User
	nextId uint
}

func NewInMemoryRepository(logger *logrus.Logger) user.Repository {
	imr := &inMemoryRepository{}
	imr.logger = logger
	imr.users = make(map[uint]user.User)
	imr.nextId = 1
	return imr
}

func (imr *inMemoryRepository) Save(ctx context.Context, user user.User) (string, error) {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("save user to in_memory_repository")
	imr.logger.Debug("check if user exists")
	exists, _ := imr.findUserByLogin(user.Login)
	if exists {
		imr.logger.Debug("user already exists in memory")
		err := uerror.ErrorDuplicate
		err.Message = fmt.Sprintf("there are user with specified login '%s'", user.Login)
		return "", err
	}
	imr.logger.Debug("generate password hash")
	if err := user.GeneratePasswordHash(); err != nil {
		imr.logger.Debugf("error during password hashing: %v", err)
		return "", err
	}
	user.Id = imr.nextId
	imr.nextId++
	imr.users[user.Id] = user
	imr.logger.Debug("user was saved")
	return user.Login, nil
}

func (imr *inMemoryRepository) GetByLogin(ctx context.Context, login string) (user.User, error) {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("get user from in_memory_repository")
	imr.logger.Debug("check if user exists")
	exists, u := imr.findUserByLogin(login)
	if exists {
		imr.logger.Debug("user found")
		return u, nil
	} else {
		imr.logger.Debugf("user was not found, login: %s", login)
		err := uerror.ErrorNotFound
		err.Message = fmt.Sprintf("user with login '%s' not found", login)
		return user.User{}, err
	}
}

func (imr *inMemoryRepository) GetById(ctx context.Context, id uint) (user.User, error) {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("get user from in_memory_repository")
	imr.logger.Debugf("find user by id: %d", id)
	u, ok := imr.users[id]
	if ok {
		imr.logger.Debug("user found")
		return u, nil
	} else {
		imr.logger.Debugf("user was not found, id: %d", id)
		err := uerror.ErrorNotFound
		err.Message = fmt.Sprintf("user with id '%d' not found", id)
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
	imr.logger.Debugf("find user by id: %d", user.Id)
	u, ok := imr.users[user.Id]
	if ok {
		imr.logger.Debug("user found")
		imr.logger.Debug("update password")
		u.Password = user.Password
		imr.logger.Debug("hashing password")
		if err := u.GeneratePasswordHash(); err != nil {
			imr.logger.Debugf("error during password hashing: %v", err)
			return err
		}
		imr.logger.Debug("update status")
		u.IsActive = user.IsActive
		imr.users[u.Id] = u
		imr.logger.Debug("user updated")
		return nil
	} else {
		imr.logger.Debugf("user was not found, id: %d", user.Id)
		err := uerror.ErrorNotFound
		err.Message = fmt.Sprintf("user with id '%d' not found", user.Id)
		return err
	}
}

func (imr *inMemoryRepository) DeleteByLogin(ctx context.Context, login string) error {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("delete user from in_memory_repository")
	imr.logger.Debugf("find user by login: %s", login)
	exists, u := imr.findUserByLogin(login)
	if exists {
		imr.logger.Debug("user found")
		delete(imr.users, u.Id)
		imr.logger.Debug("user deleted")
		return nil
	} else {
		imr.logger.Debugf("user was not found, login: %s", login)
		err := uerror.ErrorNotFound
		err.Message = fmt.Sprintf("user with login '%s' not found", login)
		return err
	}
}

func (imr *inMemoryRepository) DeleteById(ctx context.Context, id uint) error {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("delete user from in_memory_repository")
	imr.logger.Debugf("find user by id: %d", id)
	_, ok := imr.users[id]
	if ok {
		imr.logger.Debug("user found")
		delete(imr.users, id)
		imr.logger.Debug("user deleted")
		return nil
	} else {
		imr.logger.Debugf("user was not found, id: %d", id)
		err := uerror.ErrorNotFound
		err.Message = fmt.Sprintf("user with id '%d' not found", id)
		return err
	}
}

func (imr *inMemoryRepository) findUserByLogin(login string) (bool, user.User) {
	imr.logger.Debugf("check if login is free: %s", login)
	for _, v := range imr.users {
		if v.Login == login {
			return true, v
		}
	}
	return false, user.User{}
}
