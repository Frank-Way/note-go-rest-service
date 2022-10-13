package storage

import (
	"context"
	"fmt"
	"github.com/Frank-Way/note-go-rest-service/internal/note"
	"github.com/Frank-Way/note-go-rest-service/internal/note/nerror"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

var _ note.Storage = &inMemoryStorage{}

type inMemoryStorage struct {
	sync.Mutex
	logger *logrus.Logger

	notes  map[int]note.Note
	nextId int
}

func NewInMemoryStorage(logger *logrus.Logger) note.Storage {
	ims := &inMemoryStorage{}
	ims.notes = make(map[int]note.Note)
	ims.nextId = 1
	ims.logger = logger
	return ims
}

func (ims *inMemoryStorage) Save(ctx context.Context, note note.Note) (string, error) {
	ims.Lock()
	defer ims.Unlock()

	ims.logger.Info("save note to in_memory_storage")
	ims.logger.Debugf("use next id for note: %d", ims.nextId)
	note.Id = ims.nextId
	note.Author = "no-author"
	ims.notes[note.Id] = note
	ims.nextId++
	ims.logger.Debug("note was saved")
	return strconv.Itoa(int(note.Id)), nil
}

func (ims *inMemoryStorage) GetById(ctx context.Context, id int) (note.Note, error) {
	ims.Lock()
	defer ims.Unlock()

	ims.logger.Info("get note from in_memory_storage")
	ims.logger.Debugf("find note by id: %d", id)
	n, ok := ims.notes[id]
	if ok {
		ims.logger.Debug("note found")
		return n, nil
	} else {
		ims.logger.Debugf("note was not found, id: %d", id)
		err := nerror.ErrorNotFound
		err.Message = fmt.Sprintf("note with id '%d' not found", id)
		return note.Note{}, err
	}
}

func (ims *inMemoryStorage) GetAll(ctx context.Context, login string) (note.Notes, error) {
	ims.Lock()
	defer ims.Unlock()

	ims.logger.Info("get notes from in_memory_storage")
	var res []note.Note
	for _, v := range ims.notes {
		if v.Author == login {
			res = append(res, v)
		}
	}
	ims.logger.Tracef("notes: %v", res)
	ims.logger.Debug("notes found")
	return res, nil
}

func (ims *inMemoryStorage) Update(ctx context.Context, note note.Note) error {
	ims.Lock()
	defer ims.Unlock()

	ims.logger.Info("update note in in_memory_storage")
	ims.logger.Debugf("find note by id: %d", note.Id)
	n, ok := ims.notes[note.Id]
	if ok {
		ims.logger.Debug("note found")
		ims.logger.Debug("update title and text")
		n.Title = note.Title
		n.Text = note.Text
		ims.notes[n.Id] = n
		ims.logger.Debug("note updated")
		return nil
	} else {
		ims.logger.Debugf("note was not found, id: %d", n.Id)
		err := nerror.ErrorNotFound
		err.Message = fmt.Sprintf("note with id '%d' not found", n.Id)
		return err
	}
}

func (ims *inMemoryStorage) Delete(ctx context.Context, id int) error {
	ims.Lock()
	defer ims.Unlock()

	ims.logger.Info("delete note from in_memory_storage")
	ims.logger.Debugf("find note by id: %d", id)
	n, ok := ims.notes[id]
	if ok {
		ims.logger.Debug("note found")
		delete(ims.notes, n.Id)
		ims.logger.Debug("note deleted")
		return nil
	} else {
		ims.logger.Debugf("note was not found, id: %d", id)
		err := nerror.ErrorNotFound
		err.Message = fmt.Sprintf("note with id '%d' not found", id)
		return err
	}
}

func (ims *inMemoryStorage) DeleteAll(ctx context.Context) error {
	return nil
}
