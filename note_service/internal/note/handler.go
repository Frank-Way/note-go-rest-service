package note

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strconv"
)

type Handler struct {
	repository Repository
	logger     *logrus.Logger
}

func NewHandler(repository Repository, logger *logrus.Logger) *Handler {
	return &Handler{
		repository: repository,
		logger:     logger,
	}
}

var (
	noIdRe = regexp.MustCompile(`^/api/v1/notes$`)
	idRe   = regexp.MustCompile(`^/api/v1/notes/(\d+)$`)
)

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle note request")
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost && noIdRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to save handler")
		return h.SaveHandler(w, r)
	case r.Method == http.MethodGet && idRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to get handler")
		return h.GetHandler(w, r)
	case r.Method == http.MethodGet && noIdRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to get all handler")
		return h.GetAllHandler(w, r)
	case r.Method == http.MethodPut && idRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to update handler")
		return h.UpdateHandler(w, r)
	case r.Method == http.MethodDelete && idRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to delete handler")
		return h.DeleteHandler(w, r)
	case r.Method == http.MethodDelete && noIdRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to delete all handler")
		return h.DeleteAllHandler(w, r)
	default:
		h.logger.Debug("no handlers for request")
		return fmt.Errorf("wrong method %s on path: %s", r.Method, r.URL.Path)
	}
}

func (h *Handler) SaveHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle save note request")
	var n Note
	h.logger.Debug("decoding note from json")
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		h.logger.Debugf("error during decoding json: %v", err)
		return err
	}
	h.logger.Tracef("note decoded from json: %s", n)
	ctx := context.TODO()
	h.logger.Debug("pass note to repository to save it")
	uri, err := h.repository.Save(ctx, n)
	if err != nil {
		h.logger.Debugf("error during saving note to repository: %v", err)
		return err
	}
	h.logger.Tracef("no errors, note saved successfully: %s", n)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("/api/v1/notes/" + uri))
	h.logger.Debug("note saved")
	return nil
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle get note request")
	h.logger.Debug("getting id from request path")
	id, err := getIdFromUrl(r)
	if err != nil {
		h.logger.Debugf("error during getting id: %v", err)
		return err
	}
	h.logger.Tracef("got id '%s' from path '%s'", id, r.URL.Path)
	ctx := context.TODO()
	h.logger.Debug("pass id to repository to get note")
	n, err := h.repository.GetById(ctx, id)
	if err != nil {
		h.logger.Debug("error during getting note from repository")
		return err
	}
	h.logger.Tracef("got note from repository: %s", n)
	h.logger.Debug("marshaling note")
	jsonBytes, err := json.Marshal(n)
	if err != nil {
		h.logger.Debugf("error during note marshaling: %s", err)
		return err
	}
	h.logger.Trace("note marshal succeed")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
	h.logger.Debug("return note")
	return nil
}

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle get all notes request")
	ctx := context.TODO()
	h.logger.Debug("getting notes from repository")
	n, err := h.repository.GetAll(ctx)
	if err != nil {
		h.logger.Debugf("error during getting notes from repository: %v", err)
		return err
	}
	h.logger.Tracef("got notes from repository: %s", n)
	h.logger.Debug("marshaling notes")
	jsonBytes, err := json.Marshal(n)
	if err != nil {
		h.logger.Debugf("error during note marshaling: %v", err)
		return err
	}
	h.logger.Trace("notes marshal succeed")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
	h.logger.Debug("return notes")
	return nil
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle update note request")
	h.logger.Debug("getting id from request path")
	id, err := getIdFromUrl(r)
	if err != nil {
		h.logger.Debugf("error during getting id: %v", err)
		return err
	}
	h.logger.Tracef("got id '%s' from path '%s'", id, r.URL.Path)
	var n Note
	h.logger.Debug("decoding note from json")
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		h.logger.Debugf("error during decoding json: %v", err)
		return err
	}
	h.logger.Tracef("note decoded from json: %s", n)
	n.Id = id
	ctx := context.TODO()
	h.logger.Debug("pass note to repository to update it")
	if err := h.repository.Update(ctx, n); err != nil {
		h.logger.Debugf("error during updating note in repository: %v", err)
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	h.logger.Debug("note updated")
	return nil
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle delete note request")
	h.logger.Debug("getting id from request path")
	id, err := getIdFromUrl(r)
	if err != nil {
		h.logger.Debugf("error during getting id: %v", err)
		return err
	}
	h.logger.Tracef("got id '%s' from path '%s'", id, r.URL.Path)
	ctx := context.TODO()
	h.logger.Debug("pass id to repository to delete note")
	if err := h.repository.Delete(ctx, id); err != nil {
		h.logger.Debugf("error during deleting note from repository: %v", err)
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	h.logger.Debug("note deleted")
	return nil
}

func (h *Handler) DeleteAllHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Trace("handle delete all notes request")
	ctx := context.TODO()
	h.logger.Debug("deleting notes from repository")
	if err := h.repository.DeleteAll(ctx); err != nil {
		h.logger.Debugf("error during deleting notes from repository: %v", err)
		return err
	}
	return nil
}

func getIdFromUrl(r *http.Request) (uint, error) {
	matches := idRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		return 0, fmt.Errorf("no id in url")
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("id must be integer")
	}
	return uint(id), nil
}
