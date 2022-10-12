package note

import "context"

type Storage interface {
	Save(ctx context.Context, note Note) (string, error)
	GetById(ctx context.Context, id uint) (Note, error)
	GetAll(ctx context.Context, login string) (Notes, error)
	Update(ctx context.Context, note Note) error
	Delete(ctx context.Context, id uint) error
	DeleteAll(ctx context.Context) error
}
