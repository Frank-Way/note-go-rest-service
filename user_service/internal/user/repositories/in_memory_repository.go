package repositories

import (
	"context"
	"fmt"
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

	if _, ok := imr.users[user.Login]; ok {
		return "", fmt.Errorf("there are user with specified login '%s'", user.Login)
	}
	if err := user.GeneratePasswordHash(); err != nil {
		return "", err
	}
	imr.users[user.Login] = user
	return user.Login, nil
}

func (imr *inMemoryRepository) GetByLogin(ctx context.Context, login string) (user.User, error) {
	imr.Lock()
	defer imr.Unlock()

	u, ok := imr.users[login]
	if ok {
		return u, nil
	} else {
		return user.User{}, fmt.Errorf("user with login '%s' not found", login)
	}
}

func (imr *inMemoryRepository) GetAll(ctx context.Context) (user.Users, error) {
	imr.Lock()
	defer imr.Unlock()

	var res []user.User
	for _, v := range imr.users {
		res = append(res, v)
	}
	return res, nil
}

func (imr *inMemoryRepository) Update(ctx context.Context, user user.User) error {
	imr.Lock()
	defer imr.Unlock()

	u, ok := imr.users[user.Login]
	if ok {
		u.Password = user.Password
		if err := u.GeneratePasswordHash(); err != nil {
			return err
		}
		imr.users[u.Login] = u
		return nil
	} else {
		return fmt.Errorf("user with login '%s' not found", user.Login)
	}
}

func (imr *inMemoryRepository) Delete(ctx context.Context, login string) error {
	imr.Lock()
	defer imr.Unlock()

	u, ok := imr.users[login]
	if ok {
		delete(imr.users, u.Login)
		return nil
	} else {
		return fmt.Errorf("user with login '%s' not found", login)
	}
}
