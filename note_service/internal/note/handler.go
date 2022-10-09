package note

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
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

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("handle note request")
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost && NoIdRe.MatchString(r.URL.Path):
		h.logger.Trace("delegate to save handler")
		h.SaveHandler(w, r)
		return
	case r.Method == http.MethodGet && IdRe.MatchString(r.URL.Path):
		h.logger.Trace("delegate to get handler")
		h.GetByIdHandler(w, r)
		return
	case r.Method == http.MethodGet && NoIdRe.MatchString(r.URL.Path):
		h.logger.Trace("delegate to get all handler")
		h.GetAllHandler(w, r)
		return
	case r.Method == http.MethodPut && IdRe.MatchString(r.URL.Path):
		h.logger.Trace("delegate to update handler")
		h.UpdateHandler(w, r)
		return
	case r.Method == http.MethodDelete && IdRe.MatchString(r.URL.Path):
		h.logger.Trace("delegate to delete handler")
		h.DeleteHandler(w, r)
		return
	case r.Method == http.MethodDelete && NoIdRe.MatchString(r.URL.Path):
		h.logger.Trace("delegate to delete all handler")
		h.DeleteAllHandler(w, r)
		return
	default:
		h.logger.Trace("no handlers for request")
		notFound(w, r)
		return
	}
}

func (h *Handler) SaveHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("handle save note request")
	var n Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		handleError(w, r, err)
		return
	}
	ctx := context.TODO()
	uri, err := h.repository.Save(ctx, n.Title, n.Text)
	if err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("/api/v1/notes/" + uri))
}

func (h *Handler) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("handle get note request")
	id, err := getIdFromUrl(r)
	if err != nil {
		handleError(w, r, err)
		return
	}
	ctx := context.TODO()
	n, err := h.repository.GetById(ctx, id)
	if err != nil {
		handleError(w, r, err)
		return
	}

	jsonBytes, err := json.Marshal(n)
	if err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("handle get all notes request")
	ctx := context.TODO()
	n, err := h.repository.GetAll(ctx)
	if err != nil {
		handleError(w, r, err)
		return
	}

	jsonBytes, err := json.Marshal(n)
	if err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("handle update note request")
	id, err := getIdFromUrl(r)
	if err != nil {
		handleError(w, r, err)
		return
	}
	var n Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		handleError(w, r, err)
		return
	}
	ctx := context.TODO()
	if err := h.repository.Update(ctx, id, n.Text, n.Title); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("handle delete note request")
	id, err := getIdFromUrl(r)
	if err != nil {
		handleError(w, r, err)
		return
	}
	ctx := context.TODO()
	if err := h.repository.Delete(ctx, id); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteAllHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("handle delete all notes request")
	//id, err := getIdFromUrl(r)
	//if err != nil {
	//	handleError(w, r, err)
	//	return
	//}
	//ctx := context.TODO()
	//if err := h.repository.Delete(ctx, id); err != nil {
	//	handleError(w, r, err)
	//	return
	//}
	//w.WriteHeader(http.StatusOK)
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	internalServerError(w, r)
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("internal server error"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

func getIdFromUrl(r *http.Request) (uint, error) {
	matches := IdRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		return 0, fmt.Errorf("no id in url")
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("id must be integer")
	}
	return uint(id), nil
}
