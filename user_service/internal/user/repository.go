package user

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, user User) (string, error)
	GetByLogin(ctx context.Context, login string) (User, error)
	GetAll(ctx context.Context) (Users, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, login string) error
}
