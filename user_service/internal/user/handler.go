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
	//case r.Method == http.MethodGet && loginRe.MatchString(r.URL.Path):
	//	h.logger.Debug("delegate to get handler")
	//	return h.GetHandler(w, r)
	//case r.Method == http.MethodGet && noLoginRe.MatchString(r.URL.Path):
	//	h.logger.Debug("delegate to get all handler")
	//	return h.GetAllHandler(w, r)
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
	//w.Write([]byte("/api/v1/users/" + uri))
	//h.logger.Debug("pass user to repository to save it")
	//uri, err := h.repository.Save(ctx, u)
	//if err != nil {
	//	h.logger.Debugf("error during saving user to repository: %v", err)
	//	return err
	//}
	//h.logger.Tracef("no errors, user saved successfully: %s", u)
	//w.WriteHeader(http.StatusCreated)
	//w.Write([]byte("/api/v1/users/" + uri))
	//h.logger.Debug("user saved")
	//return nil
}

//func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) error {
//	h.logger.Info("handle get user request")
//	h.logger.Debug("getting login from request path")
//	login, err := getLoginFromUrl(r)
//	if err != nil {
//		h.logger.Debugf("error during getting login: %v", err)
//		return err
//	}
//	h.logger.Tracef("got login '%s' from path '%s'", login, r.URL.Path)
//	h.logger.Debug("pass login to repository to get user")
//	u, err := h.repository.GetByLogin(ctx, login)
//	if err != nil {
//		h.logger.Debug("error during getting user from repository")
//		return err
//	}
//	h.logger.Tracef("got login from repository: %s", u.Login)
//	h.logger.Debug("marshaling user")
//	jsonBytes, err := json.Marshal(u)
//	if err != nil {
//		h.logger.Debugf("error during user marshaling: %s", err)
//		return err
//	}
//	h.logger.Trace("user marshal succeed")
//	w.WriteHeader(http.StatusOK)
//	w.Write(jsonBytes)
//	h.logger.Debug("return user")
//	return nil
//}

//func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) error {
//	h.logger.Info("handle get all users request")
//	h.logger.Debug("getting users from repository")
//	u, err := h.repository.GetAll(ctx)
//	if err != nil {
//		h.logger.Debugf("error during getting users from repository: %v", err)
//		return err
//	}
//	h.logger.Tracef("got users from repository")
//	h.logger.Debug("marshaling users")
//	jsonBytes, err := json.Marshal(u)
//	if err != nil {
//		h.logger.Debugf("error during user marshaling: %v", err)
//		return err
//	}
//	h.logger.Trace("users marshal succeed")
//	w.WriteHeader(http.StatusOK)
//	w.Write(jsonBytes)
//	h.logger.Debug("return users")
//	return nil
//}

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
	//h.logger.Debug("getting login from request path")
	//login, err := getLoginFromUrl(r)
	//if err != nil {
	//	h.logger.Debugf("error during getting login: %v", err)
	//	return err
	//}
	//h.logger.Tracef("got login '%s' from path '%s'", login, r.URL.Path)
	//var u User
	//h.logger.Debug("decoding user from json")
	//if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
	//	h.logger.Debugf("error during decoding json: %v", err)
	//	return err
	//}
	//h.logger.Tracef("user decoded from json: %s", u)
	//u.Login = login
	//h.logger.Debug("pass user to repository to update it")
	//if err := h.repository.Update(ctx, u); err != nil {
	//	h.logger.Debugf("error during updating user in repository: %v", err)
	//	return err
	//}
	//w.WriteHeader(http.StatusNoContent)
	//h.logger.Debug("user updated")
	//return nil
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
	h.logger.Debug("pass login to service")
	if err = h.service.DeleteUser(r.Context(), login); err != nil {
		h.logger.Debugf("error in service: %v", err)
		return err
	}
	//h.logger.Debug("pass login to repository to delete user")
	//if err := h.repository.Delete(ctx, login); err != nil {
	//	h.logger.Debugf("error during deleting user from repository: %v", err)
	//	return err
	//}
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
	if err = h.service.SignIn(r.Context(), login, uDTO); err != nil {
		h.logger.Debugf("error in service: %v", err)
		return err
	}
	w.WriteHeader(http.StatusOK)
	// TODO write auth headers???
	return nil
}

func getLoginFromUrl(r *http.Request) (string, error) {
	matches := loginRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		return "", fmt.Errorf("no login in url: %s", r.URL.Path)
	}
	return matches[1], nil
}
