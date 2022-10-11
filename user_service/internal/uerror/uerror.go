package uerror

import "encoding/json"

var (
	ErrorNotFound   = NewUserError(nil, "user not found", "", "US-1")
	ErrorDuplicate  = NewUserError(nil, "user already exists", "", "US-2")
	ErrorRepository = NewUserError(nil, "repository error", "", "US-3")
	ErrorNoAuth     = NewUserError(nil, "no authorized", "", "US-4")
)

type UserError struct {
	Err              error  `json:"-"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developer_message"`
	Code             string `json:"code"`
}

func (e *UserError) Error() string {
	return e.Message
}

func (e *UserError) Unwrap() error { return e.Err }

func (e *UserError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewUserError(err error, message, developerMessage, code string) *UserError {
	return &UserError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}

func systemError(err error) *UserError {
	return NewUserError(err, "internal system error", err.Error(), "US-0")
}
