package user

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, user User) (string, error)
	GetByLogin(ctx context.Context, login string) (User, error)
	GetById(ctx context.Context, id uint) (User, error)
	GetAll(ctx context.Context) (Users, error)
	Update(ctx context.Context, user User) error
	DeleteByLogin(ctx context.Context, login string) error
	DeleteById(ctx context.Context, id uint) error
}
