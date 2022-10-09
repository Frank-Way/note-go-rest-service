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

func (imr *inMemoryRepository) Save(ctx context.Context, title, text string) (string, error) {
	imr.Lock()
	defer imr.Unlock()

	n := note.Note{
		Id:     imr.nextId,
		Title:  title,
		Text:   text,
		Author: "no-author",
	}
	imr.notes[n.Id] = n
	imr.nextId++
	return strconv.Itoa(int(n.Id)), nil
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

func (imr *inMemoryRepository) Update(ctx context.Context, id uint, title, text string) error {
	imr.Lock()
	defer imr.Unlock()

	n, ok := imr.notes[id]
	if ok {
		n.Title = title
		n.Text = text
		imr.notes[n.Id] = n
		return nil
	} else {
		return fmt.Errorf("note with id '%d' not found", id)
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
