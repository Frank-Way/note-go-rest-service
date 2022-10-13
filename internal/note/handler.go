package note

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strconv"
)

type Handler struct {
	service Service
	logger  *logrus.Logger
}

func NewHandler(service Service, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
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
		return h.saveHandler(w, r)
	case r.Method == http.MethodGet && idRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to get handler")
		return h.getHandler(w, r)
	case r.Method == http.MethodGet && noIdRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to get all handler")
		return h.getAllHandler(w, r)
	case r.Method == http.MethodPut && idRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to update handler")
		return h.updateHandler(w, r)
	case r.Method == http.MethodDelete && idRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to delete handler")
		return h.deleteHandler(w, r)
	default:
		h.logger.Debug("no handlers for request")
		return fmt.Errorf("wrong method %s on path: %s", r.Method, r.URL.Path)
	}
}

func (h *Handler) saveHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle save note request")
	var dto CreateNoteDTO
	h.logger.Debug("decoding note dto from json")
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Debugf("error during decoding json: %v", err)
		return err
	}
	h.logger.Tracef("note dto decoded from json: %v", dto)
	h.logger.Debug("get header 'Authorization' from request")
	authHeader := r.Header.Get("Authorization")
	h.logger.Debug("pass auth and note dto to service to save it")
	uri, err := h.service.CreateNote(r.Context(), authHeader, dto)
	if err != nil {
		h.logger.Debugf("error during saving note in service: %v", err)
		return err
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", "/api/v1/notes/"+uri)
	w.Write([]byte("/api/v1/notes/" + uri))
	h.logger.Debug("note saved")
	return nil
}

func (h *Handler) getHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle get note request")
	h.logger.Debug("getting id from request path")
	id, err := getIdFromUrl(r)
	if err != nil {
		h.logger.Debugf("error during getting id: %v", err)
		return err
	}
	h.logger.Tracef("got id '%d' from path '%s'", id, r.URL.Path)
	h.logger.Debug("get header 'Authorization' from request")
	authHeader := r.Header.Get("Authorization")
	h.logger.Debug("pass auth and id to service to get note")
	n, err := h.service.GetNote(r.Context(), authHeader, id)
	if err != nil {
		h.logger.Debugf("error during getting note from service: %v", err)
		return err
	}
	h.logger.Tracef("got note from service: %v", n)
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

func (h *Handler) getAllHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle get all notes request")
	h.logger.Debug("get header 'Authorization' from request")
	authHeader := r.Header.Get("Authorization")
	h.logger.Debug("pass auth to service to get all notes")
	n, err := h.service.GetAllNotes(r.Context(), authHeader)
	if err != nil {
		h.logger.Debugf("error during getting notes from service: %v", err)
		return err
	}
	h.logger.Tracef("got notes from storage: %v", n)
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

func (h *Handler) updateHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle update note request")
	h.logger.Debug("getting id from request path")
	id, err := getIdFromUrl(r)
	if err != nil {
		h.logger.Debugf("error during getting id: %v", err)
		return err
	}
	h.logger.Tracef("got id '%d' from path '%s'", id, r.URL.Path)
	var dto UpdateNoteDTO
	h.logger.Debug("decoding note dto from json")
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.logger.Debugf("error during decoding json: %v", err)
		return err
	}
	h.logger.Tracef("note dto decoded from json: %v", dto)
	h.logger.Debug("get header 'Authorization' from request")
	authHeader := r.Header.Get("Authorization")
	h.logger.Debug("pass auth and note dto to service to update it")
	if err := h.service.UpdateNote(r.Context(), authHeader, id, dto); err != nil {
		h.logger.Debugf("error during updating note in service: %v", err)
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	h.logger.Debug("note updated")
	return nil
}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle delete note request")
	h.logger.Debug("getting id from request path")
	id, err := getIdFromUrl(r)
	if err != nil {
		h.logger.Debugf("error during getting id: %v", err)
		return err
	}
	h.logger.Tracef("got id '%d' from path '%s'", id, r.URL.Path)
	h.logger.Debug("get header 'Authorization' from request")
	authHeader := r.Header.Get("Authorization")
	h.logger.Debug("pass auth and id to service to delete note")
	if err := h.service.DeleteNote(r.Context(), authHeader, id); err != nil {
		h.logger.Debugf("error during deleting note from service: %v", err)
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	h.logger.Debug("note deleted")
	return nil
}

func getIdFromUrl(r *http.Request) (int, error) {
	matches := idRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		return 0, fmt.Errorf("no id in url")
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("id must be integer")
	}
	return int(id), nil
}
