package note

import (
	"context"
	"regexp"
)

var (
	NoIdRe = regexp.MustCompile(`^/api/v1/notes$`)
	IdRe   = regexp.MustCompile(`^/api/v1/notes/(\d+)$`)
)

type Repository interface {
	Save(ctx context.Context, title, text string) (string, error)
	GetById(ctx context.Context, id uint) (Note, error)
	GetAll(ctx context.Context) (Notes, error)
	Update(ctx context.Context, id uint, title, text string) error
	Delete(ctx context.Context, id uint) error
	DeleteAll(ctx context.Context) error
}
