package user

import (
	"context"
	"regexp"
)

var (
	NoLoginRe = regexp.MustCompile(`^/api/v1/users$`)
	LoginRe   = regexp.MustCompile(`^/api/v1/users/([A-Za-z0-9_]+)$`)
)

type Repository interface {
	Save(ctx context.Context, login, password string) (string, error)
	GetByLogin(ctx context.Context, login string) (User, error)
	GetAll(ctx context.Context) (Users, error)
	Update(ctx context.Context, login, password string) error
	Delete(ctx context.Context, login string) error
}
