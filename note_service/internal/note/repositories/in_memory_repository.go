package repositories

import (
	"context"
	"fmt"
	"github.com/Frank-Way/note-go-rest-service/note_service/internal/note"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

var _ note.Repository = &inMemoryRepository{}

type inMemoryRepository struct {
	sync.Mutex
	logger *logrus.Logger

	notes  map[uint]note.Note
	nextId uint
}

func NewInMemoryRepository(logger *logrus.Logger) note.Repository {
	imr := &inMemoryRepository{}
	imr.notes = make(map[uint]note.Note)
	imr.nextId = 1
	imr.logger = logger
	return imr
}

func (imr *inMemoryRepository) Save(ctx context.Context, note note.Note) (string, error) {
	imr.Lock()
	defer imr.Unlock()

	note.Id = imr.nextId
	note.Author = "no-author"
	imr.notes[note.Id] = note
	imr.nextId++
	return strconv.Itoa(int(note.Id)), nil
}

func (imr *inMemoryRepository) GetById(ctx context.Context, id uint) (note.Note, error) {
	imr.Lock()
	defer imr.Unlock()

	n, ok := imr.notes[id]
	if ok {
		return n, nil
	} else {
		return note.Note{}, fmt.Errorf("note with id '%d' not found", id)
	}
}

func (imr *inMemoryRepository) GetAll(ctx context.Context) (note.Notes, error) {
	imr.Lock()
	defer imr.Unlock()

	var res []note.Note
	for _, v := range imr.notes {
		res = append(res, v)
	}
	return res, nil
}

func (imr *inMemoryRepository) Update(ctx context.Context, note note.Note) error {
	imr.Lock()
	defer imr.Unlock()

	n, ok := imr.notes[note.Id]
	if ok {
		n.Title = note.Title
		n.Text = note.Text
		imr.notes[n.Id] = n
		return nil
	} else {
		return fmt.Errorf("note with id '%d' not found", note.Id)
	}
}

func (imr *inMemoryRepository) Delete(ctx context.Context, id uint) error {
	imr.Lock()
	defer imr.Unlock()

	n, ok := imr.notes[id]
	if ok {
		delete(imr.notes, n.Id)
		return nil
	} else {
		return fmt.Errorf("note with id '%d' not found", id)
	}
}

func (imr *inMemoryRepository) DeleteAll(ctx context.Context) error {
	return nil
}
