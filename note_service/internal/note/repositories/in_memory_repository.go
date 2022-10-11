package repositories

import (
	"context"
	"fmt"
	"github.com/Frank-Way/note-go-rest-service/note_service/internal/nerror"
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

	imr.logger.Info("save note to in_memory_repository")
	// TODO IMPLEMENT GETTING USER FROM CONTEXT OR SMTH IDK
	imr.logger.Warn("TODO IMPLEMENT GETTING USER FROM CONTEXT OR SMTH IDK")
	imr.logger.Debug("getting user")
	imr.logger.Debug("check user's permissions")
	imr.logger.Debugf("use next id for note: %d", imr.nextId)
	note.Id = imr.nextId
	note.Author = "no-author"
	imr.notes[note.Id] = note
	imr.nextId++
	imr.logger.Debug("note was saved")
	return strconv.Itoa(int(note.Id)), nil
}

func (imr *inMemoryRepository) GetById(ctx context.Context, id uint) (note.Note, error) {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("get note from in_memory_repository")
	// TODO IMPLEMENT GETTING USER FROM CONTEXT OR SMTH IDK
	imr.logger.Warn("TODO IMPLEMENT GETTING USER FROM CONTEXT OR SMTH IDK")
	imr.logger.Debug("getting user")
	imr.logger.Debug("check user's permissions")
	imr.logger.Debugf("find note by id: %d", id)
	n, ok := imr.notes[id]
	if ok {
		imr.logger.Debug("note found")
		return n, nil
	} else {
		imr.logger.Debugf("note was not found, id: %d", id)
		err := nerror.ErrorNotFound
		err.Message = fmt.Sprintf("note with id '%d' not found", id)
		return note.Note{}, err
	}
}

func (imr *inMemoryRepository) GetAll(ctx context.Context) (note.Notes, error) {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("get notes from in_memory_repository")
	// TODO IMPLEMENT GETTING USER FROM CONTEXT OR SMTH IDK
	imr.logger.Warn("TODO IMPLEMENT GETTING USER FROM CONTEXT OR SMTH IDK")
	imr.logger.Debug("getting user")
	imr.logger.Debug("check user's permissions")
	var res []note.Note
	for _, v := range imr.notes {
		res = append(res, v)
	}
	imr.logger.Tracef("notes: %v", res)
	imr.logger.Debug("notes found")
	return res, nil
}

func (imr *inMemoryRepository) Update(ctx context.Context, note note.Note) error {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("update note in in_memory_repository")
	// TODO IMPLEMENT GETTING USER FROM CONTEXT OR SMTH IDK
	imr.logger.Warn("TODO IMPLEMENT GETTING USER FROM CONTEXT OR SMTH IDK")
	imr.logger.Debug("getting user")
	imr.logger.Debug("check user's permissions")
	imr.logger.Debugf("find note by id: %d", note.Id)
	n, ok := imr.notes[note.Id]
	if ok {
		imr.logger.Debug("note found")
		imr.logger.Debug("update title and text")
		n.Title = note.Title
		n.Text = note.Text
		imr.notes[n.Id] = n
		imr.logger.Debug("note updated")
		return nil
	} else {
		imr.logger.Debugf("note was not found, id: %d", n.Id)
		err := nerror.ErrorNotFound
		err.Message = fmt.Sprintf("note with id '%d' not found", n.Id)
		return err
	}
}

func (imr *inMemoryRepository) Delete(ctx context.Context, id uint) error {
	imr.Lock()
	defer imr.Unlock()

	imr.logger.Info("delete note from in_memory_repository")
	// TODO IMPLEMENT GETTING USER FROM CONTEXT OR SMTH IDK
	imr.logger.Warn("TODO IMPLEMENT GETTING USER FROM CONTEXT OR SMTH IDK")
	imr.logger.Debug("getting user")
	imr.logger.Debug("check user's permissions")
	imr.logger.Debugf("find note by id: %d", id)
	n, ok := imr.notes[id]
	if ok {
		imr.logger.Debug("note found")
		delete(imr.notes, n.Id)
		imr.logger.Debug("note deleted")
		return nil
	} else {
		imr.logger.Debugf("note was not found, id: %d", id)
		err := nerror.ErrorNotFound
		err.Message = fmt.Sprintf("note with id '%d' not found", id)
		return err
	}
}

func (imr *inMemoryRepository) DeleteAll(ctx context.Context) error {
	return nil
}
