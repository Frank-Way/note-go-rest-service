package nerror

import (
	"errors"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var userError *NoteError
		err := h(w, r)
		if err != nil {
			if errors.As(err, &userError) {
				if errors.Is(err, ErrorNotFound) {
					w.WriteHeader(http.StatusNotFound)
				} else if errors.Is(err, ErrorDuplicate) {
					w.WriteHeader(http.StatusForbidden)
				} else if errors.Is(err, ErrorStorage) {
					w.WriteHeader(http.StatusInternalServerError)
				} else if errors.Is(err, ErrorNoAuth) {
					w.WriteHeader(http.StatusUnauthorized)
				}

				err = err.(*NoteError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(userError.Marshal())
				return
			}

			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err).Marshal())
		}
	}
}
