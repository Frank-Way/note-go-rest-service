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

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle user request")
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost && noLoginRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to save handler")
		return h.SaveHandler(w, r)
	case r.Method == http.MethodGet && loginRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to get handler")
		return h.GetHandler(w, r)
	case r.Method == http.MethodGet && noLoginRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to get all handler")
		return h.GetAllHandler(w, r)
	case r.Method == http.MethodPut && loginRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to update handler")
		return h.UpdateHandler(w, r)
	case r.Method == http.MethodDelete && loginRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to delete handler")
		return h.DeleteHandler(w, r)
	default:
		h.logger.Debug("no handlers for request")
		return fmt.Errorf("wrong method %s on path: %s", r.Method, r.URL.Path)
	}
}

func (h *Handler) SaveHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle save user request")
	var u User
	h.logger.Debug("decoding user from json")
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		h.logger.Debugf("error during decoding json: %v", err)
		return err
	}
	h.logger.Tracef("user decoded from json: %s", u)
	ctx := context.TODO()
	h.logger.Debug("pass user to repository to save it")
	uri, err := h.repository.Save(ctx, u)
	if err != nil {
		h.logger.Debugf("error during saving user to repository: %v", err)
		return err
	}
	h.logger.Tracef("no errors, user saved successfully: %s", u)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("/api/v1/users/" + uri))
	h.logger.Debug("user saved")
	return nil
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle get user request")
	h.logger.Debug("getting login from request path")
	login, err := getLoginFromUrl(r)
	if err != nil {
		h.logger.Debugf("error during getting login: %v", err)
		return err
	}
	h.logger.Tracef("got login '%s' from path '%s'", login, r.URL.Path)
	ctx := context.TODO()
	h.logger.Debug("pass login to repository to get user")
	u, err := h.repository.GetByLogin(ctx, login)
	if err != nil {
		h.logger.Debug("error during getting user from repository")
		return err
	}
	h.logger.Tracef("got login from repository: %s", u.Login)
	h.logger.Debug("marshaling user")
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		h.logger.Debugf("error during user marshaling: %s", err)
		return err
	}
	h.logger.Trace("user marshal succeed")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
	h.logger.Debug("return user")
	return nil
}

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle get all users request")
	ctx := context.TODO()
	h.logger.Debug("getting users from repository")
	u, err := h.repository.GetAll(ctx)
	if err != nil {
		h.logger.Debugf("error during getting users from repository: %v", err)
		return err
	}
	h.logger.Tracef("got users from repository")
	h.logger.Debug("marshaling users")
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		h.logger.Debugf("error during user marshaling: %v", err)
		return err
	}
	h.logger.Trace("users marshal succeed")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
	h.logger.Debug("return users")
	return nil
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle update user request")
	h.logger.Debug("getting login from request path")
	login, err := getLoginFromUrl(r)
	if err != nil {
		h.logger.Debugf("error during getting login: %v", err)
		return err
	}
	h.logger.Tracef("got login '%s' from path '%s'", login, r.URL.Path)
	var u User
	h.logger.Debug("decoding user from json")
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		h.logger.Debugf("error during decoding json: %v", err)
		return err
	}
	h.logger.Tracef("user decoded from json: %s", u)
	u.Login = login
	ctx := context.TODO()
	h.logger.Debug("pass user to repository to update it")
	if err := h.repository.Update(ctx, u); err != nil {
		h.logger.Debugf("error during updating user in repository: %v", err)
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	h.logger.Debug("user updated")
	return nil
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle delete user request")
	h.logger.Debug("getting login from request path")
	login, err := getLoginFromUrl(r)
	if err != nil {
		h.logger.Debugf("error during getting login: %v", err)
		return err
	}
	h.logger.Tracef("got login '%s' from path '%s'", login, r.URL.Path)
	ctx := context.TODO()
	h.logger.Debug("pass login to repository to delete user")
	if err := h.repository.Delete(ctx, login); err != nil {
		h.logger.Debugf("error during deleting user from repository: %v", err)
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	h.logger.Debug("user deleted")
	return nil
}

func getLoginFromUrl(r *http.Request) (string, error) {
	matches := loginRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		return "", fmt.Errorf("no login in url: %s", r.URL.Path)
	}
	return matches[1], nil
}
