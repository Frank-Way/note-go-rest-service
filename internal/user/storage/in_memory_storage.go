package storage

import (
	"context"
	"fmt"
	"github.com/Frank-Way/note-go-rest-service/internal/user"
	"github.com/Frank-Way/note-go-rest-service/internal/user/uerror"
	"github.com/sirupsen/logrus"
	"sync"
)

var _ user.Storage = &inMemoryStorage{}

type inMemoryStorage struct {
	sync.Mutex
	logger *logrus.Logger

	users  map[int]user.User
	nextId int
}

func NewInMemoryStorage(logger *logrus.Logger) user.Storage {
	ims := &inMemoryStorage{}
	ims.logger = logger
	ims.users = make(map[int]user.User)
	ims.nextId = 1
	return ims
}

func (ims *inMemoryStorage) Save(ctx context.Context, user user.User) (string, error) {
	ims.Lock()
	defer ims.Unlock()

	ims.logger.Info("save user to in_memory_storage")
	ims.logger.Debug("check if user exists")
	exists, _ := ims.findUserByLogin(user.Login)
	if exists {
		ims.logger.Debug("user already exists in memory")
		err := uerror.ErrorDuplicate
		err.Message = fmt.Sprintf("there are user with specified login '%s'", user.Login)
		return "", err
	}
	ims.logger.Debug("generate password hash")
	if err := user.GeneratePasswordHash(); err != nil {
		ims.logger.Debugf("error during password hashing: %v", err)
		return "", err
	}
	user.Id = ims.nextId
	ims.nextId++
	ims.users[user.Id] = user
	ims.logger.Debug("user was saved")
	return user.Login, nil
}

func (ims *inMemoryStorage) GetByLogin(ctx context.Context, login string) (user.User, error) {
	ims.Lock()
	defer ims.Unlock()

	ims.logger.Info("get user from in_memory_storage")
	ims.logger.Debug("check if user exists")
	exists, u := ims.findUserByLogin(login)
	if exists {
		ims.logger.Debug("user found")
		return u, nil
	} else {
		ims.logger.Debugf("user was not found, login: %s", login)
		err := uerror.ErrorNotFound
		err.Message = fmt.Sprintf("user with login '%s' not found", login)
		return user.User{}, err
	}
}

func (ims *inMemoryStorage) GetById(ctx context.Context, id int) (user.User, error) {
	ims.Lock()
	defer ims.Unlock()

	ims.logger.Info("get user from in_memory_storage")
	ims.logger.Debugf("find user by id: %d", id)
	u, ok := ims.users[id]
	if ok {
		ims.logger.Debug("user found")
		return u, nil
	} else {
		ims.logger.Debugf("user was not found, id: %d", id)
		err := uerror.ErrorNotFound
		err.Message = fmt.Sprintf("user with id '%d' not found", id)
		return user.User{}, err
	}
}

func (ims *inMemoryStorage) GetAll(ctx context.Context) (user.Users, error) {
	ims.Lock()
	defer ims.Unlock()

	ims.logger.Info("get users from in_memory_storage")
	var res []user.User
	for _, v := range ims.users {
		res = append(res, v)
	}
	ims.logger.Debug("users found")
	return res, nil
}

func (ims *inMemoryStorage) Update(ctx context.Context, user user.User) error {
	ims.Lock()
	defer ims.Unlock()

	ims.logger.Info("update user in in_memory_storage")
	ims.logger.Debugf("find user by id: %d", user.Id)
	u, ok := ims.users[user.Id]
	if ok {
		ims.logger.Debug("user found")
		ims.logger.Debug("update password")
		u.Password = user.Password
		ims.logger.Debug("hashing password")
		if err := u.GeneratePasswordHash(); err != nil {
			ims.logger.Debugf("error during password hashing: %v", err)
			return err
		}
		ims.logger.Debug("update status")
		u.IsActive = user.IsActive
		ims.users[u.Id] = u
		ims.logger.Debug("user updated")
		return nil
	} else {
		ims.logger.Debugf("user was not found, id: %d", user.Id)
		err := uerror.ErrorNotFound
		err.Message = fmt.Sprintf("user with id '%d' not found", user.Id)
		return err
	}
}

func (ims *inMemoryStorage) DeleteByLogin(ctx context.Context, login string) error {
	ims.Lock()
	defer ims.Unlock()

	ims.logger.Info("delete user from in_memory_storage")
	ims.logger.Debugf("find user by login: %s", login)
	exists, u := ims.findUserByLogin(login)
	if exists {
		ims.logger.Debug("user found")
		delete(ims.users, u.Id)
		ims.logger.Debug("user deleted")
		return nil
	} else {
		ims.logger.Debugf("user was not found, login: %s", login)
		err := uerror.ErrorNotFound
		err.Message = fmt.Sprintf("user with login '%s' not found", login)
		return err
	}
}

func (ims *inMemoryStorage) DeleteById(ctx context.Context, id int) error {
	ims.Lock()
	defer ims.Unlock()

	ims.logger.Info("delete user from in_memory_storage")
	ims.logger.Debugf("find user by id: %d", id)
	_, ok := ims.users[id]
	if ok {
		ims.logger.Debug("user found")
		delete(ims.users, id)
		ims.logger.Debug("user deleted")
		return nil
	} else {
		ims.logger.Debugf("user was not found, id: %d", id)
		err := uerror.ErrorNotFound
		err.Message = fmt.Sprintf("user with id '%d' not found", id)
		return err
	}
}

func (ims *inMemoryStorage) findUserByLogin(login string) (bool, user.User) {
	ims.logger.Debugf("check if login is free: %s", login)
	for _, v := range ims.users {
		if v.Login == login {
			return true, v
		}
	}
	return false, user.User{}
}
