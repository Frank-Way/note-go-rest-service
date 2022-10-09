package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
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
	noLoginRe = regexp.MustCompile(`^/api/v1/users$`)
	loginRe   = regexp.MustCompile(`^/api/v1/users/([A-Za-z0-9_]+)$`)
)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("handle user request")
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost && noLoginRe.MatchString(r.URL.Path):
		h.logger.Trace("delegate to save handler")
		h.SaveHandler(w, r)
		return
	case r.Method == http.MethodGet && loginRe.MatchString(r.URL.Path):
		h.logger.Trace("delegate to get handler")
		h.GetHandler(w, r)
		return
	case r.Method == http.MethodGet && noLoginRe.MatchString(r.URL.Path):
		h.logger.Trace("delegate to get all handler")
		h.GetAllHandler(w, r)
		return
	case r.Method == http.MethodPut && loginRe.MatchString(r.URL.Path):
		h.logger.Trace("delegate to update handler")
		h.UpdateHandler(w, r)
		return
	case r.Method == http.MethodDelete && loginRe.MatchString(r.URL.Path):
		h.logger.Trace("delegate to delete handler")
		h.DeleteHandler(w, r)
		return
	default:
		h.logger.Trace("no handlers for request")
		notFound(w, r)
		return
	}
}

func (h *Handler) SaveHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("handle save user request")
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		handleError(w, r, err)
		return
	}
	ctx := context.TODO()
	uri, err := h.repository.Save(ctx, u)
	if err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("/api/v1/users/" + uri))
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("handle get user request")
	login, err := getLoginFromUrl(r)
	if err != nil {
		handleError(w, r, err)
		return
	}
	ctx := context.TODO()
	u, err := h.repository.GetByLogin(ctx, login)
	if err != nil {
		handleError(w, r, err)
		return
	}

	jsonBytes, err := json.Marshal(u)
	if err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("handle get all users request")
	ctx := context.TODO()
	u, err := h.repository.GetAll(ctx)
	if err != nil {
		handleError(w, r, err)
		return
	}

	jsonBytes, err := json.Marshal(u)
	if err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("handle update user request")
	login, err := getLoginFromUrl(r)
	if err != nil {
		handleError(w, r, err)
		return
	}
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		handleError(w, r, err)
		return
	}
	u.Login = login
	ctx := context.TODO()
	if err := h.repository.Update(ctx, u); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Trace("handle delete user request")
	login, err := getLoginFromUrl(r)
	if err != nil {
		handleError(w, r, err)
		return
	}
	ctx := context.TODO()
	if err := h.repository.Delete(ctx, login); err != nil {
		handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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

func getLoginFromUrl(r *http.Request) (string, error) {
	matches := loginRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		return "", fmt.Errorf("no login in url")
	}
	return matches[1], nil
}
