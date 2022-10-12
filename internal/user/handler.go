package user

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
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
	noLoginRe = regexp.MustCompile(`^/api/v1/users$`)
	loginRe   = regexp.MustCompile(`^/api/v1/users/([A-Za-z0-9_]+)$`)
)

func (h *Handler) Handler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle user request")
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost && noLoginRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to save handler")
		return h.saveHandler(w, r)
	case r.Method == http.MethodPut && loginRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to update handler")
		return h.updateHandler(w, r)
	case r.Method == http.MethodPost && loginRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to auth handler")
		return h.authHandler(w, r)
	case r.Method == http.MethodDelete && loginRe.MatchString(r.URL.Path):
		h.logger.Debug("delegate to delete handler")
		return h.deleteHandler(w, r)
	// TODO DELETE DEBUG ENDPOINTS
	//case r.Method == http.MethodGet && r.URL.Path == `/debug/allusers`:
	//	u, err := h.service.TMPGetAllUsers(r.Context())
	//	if err != nil {
	//		return err
	//	}
	//	jsonBytes, err := json.Marshal(u)
	//	if err != nil {
	//		return err
	//	}
	//	w.WriteHeader(http.StatusOK)
	//	w.Write(jsonBytes)
	//	return nil
	default:
		h.logger.Debug("no handlers for request")
		return fmt.Errorf("wrong method %s on path: %s", r.Method, r.URL.Path)
	}
}

func (h *Handler) saveHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle save user request")
	var uDTO CreateUserDTO
	h.logger.Debug("decoding create user dto from json")
	if err := json.NewDecoder(r.Body).Decode(&uDTO); err != nil {
		h.logger.Debugf("error during decoding json: %v", err)
		return err
	}
	h.logger.Debug("pass dto to service")
	uri, err := h.service.SignUp(r.Context(), uDTO)
	if err != nil {
		h.logger.Debugf("error in service: %v", err)
		return err
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", "/api/v1/users/"+uri)
	return nil
}

func (h *Handler) updateHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle update user request")
	h.logger.Debug("getting login from request path")
	login, err := getLoginFromUrl(r)
	if err != nil {
		h.logger.Debugf("error during getting login: %v", err)
		return err
	}
	h.logger.Tracef("got login '%s' from path '%s'", login, r.URL.Path)
	var uDTO UpdateUserDTO
	h.logger.Debug("decoding update user dto from json")
	if err := json.NewDecoder(r.Body).Decode(&uDTO); err != nil {
		h.logger.Debugf("error during decoding json: %v", err)
		return err
	}
	h.logger.Debug("pass dto to service")
	if err := h.service.ChangePassword(r.Context(), login, uDTO); err != nil {
		h.logger.Debugf("error in service: %v", err)
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle delete user request")
	h.logger.Debug("getting login from request path")
	login, err := getLoginFromUrl(r)
	if err != nil {
		h.logger.Debugf("error during getting login: %v", err)
		return err
	}
	h.logger.Tracef("got login '%s' from path '%s'", login, r.URL.Path)
	authHeader := r.Header.Get("Authorization")
	h.logger.Debug("pass login to service")
	if err = h.service.DeleteUser(r.Context(), authHeader, login); err != nil {
		h.logger.Debugf("error in service: %v", err)
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) authHandler(w http.ResponseWriter, r *http.Request) error {
	h.logger.Info("handle auth user request")
	h.logger.Debug("getting login from request path")
	login, err := getLoginFromUrl(r)
	if err != nil {
		h.logger.Debugf("error during getting login: %v", err)
		return err
	}
	h.logger.Tracef("got login '%s' from path '%s'", login, r.URL.Path)
	var uDTO AuthUserDTO
	h.logger.Debug("decoding auth user dto from json")
	if err := json.NewDecoder(r.Body).Decode(&uDTO); err != nil {
		h.logger.Debugf("error during decoding json: %v", err)
		return err
	}
	h.logger.Debug("pass dto to service")
	token, err := h.service.SignIn(r.Context(), login, uDTO)
	if err != nil {
		h.logger.Debugf("error in service: %v", err)
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
	return nil
}

func getLoginFromUrl(r *http.Request) (string, error) {
	matches := loginRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		return "", fmt.Errorf("no login in url: %s", r.URL.Path)
	}
	return matches[1], nil
}
