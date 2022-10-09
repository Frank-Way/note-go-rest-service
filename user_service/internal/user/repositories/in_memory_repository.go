package repositories

import (
	"context"
	"fmt"
	"github.com/Frank-Way/note-go-rest-service/user_service/internal/user"
	"sync"
)

var _ user.Repository = &inMemoryRepository{}

type inMemoryRepository struct {
	sync.Mutex

	users map[string]user.User
}

func NewInMemoryRepository() user.Repository {
	imr := &inMemoryRepository{}
	imr.users = make(map[string]user.User)
	return imr
}

func (imr *inMemoryRepository) Save(ctx context.Context, login, password string) (string, error) {
	imr.Lock()
	defer imr.Unlock()

	if _, ok := imr.users[login]; ok {
		return "", fmt.Errorf("there are user with specified login '%s'", login)
	}

	u := user.User{
		Login:    login,
		Password: password,
	}
	if err := u.GeneratePasswordHash(); err != nil {
		return "", err
	}
	imr.users[u.Login] = u
	return u.Login, nil
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

func (imr *inMemoryRepository) Update(ctx context.Context, login, password string) error {
	imr.Lock()
	defer imr.Unlock()

	u, ok := imr.users[login]
	if ok {
		u.Password = password
		if err := u.GeneratePasswordHash(); err != nil {
			return err
		}
		imr.users[u.Login] = u
		return nil
	} else {
		return fmt.Errorf("user with login '%s' not found", login)
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
