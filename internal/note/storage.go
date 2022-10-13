package note

import "context"

type Storage interface {
	Save(ctx context.Context, note Note) (string, error)
	GetById(ctx context.Context, id int) (Note, error)
	GetAll(ctx context.Context, login string) (Notes, error)
	Update(ctx context.Context, note Note) error
	Delete(ctx context.Context, id int) error
	DeleteAll(ctx context.Context) error
}
