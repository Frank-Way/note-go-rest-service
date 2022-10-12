package uerror

import "encoding/json"

var (
	ErrorNotFound          = NewUserError(nil, "user not found", "", "US-1")
	ErrorDuplicate         = NewUserError(nil, "user already exists", "", "US-2")
	ErrorStorage           = NewUserError(nil, "storage error", "", "US-3")
	ErrorNoAuth            = NewUserError(nil, "no authorized", "", "US-4")
	ErrorPasswordsMismatch = NewUserError(nil, "passwords mismatches", "", "US-5")
	ErrorWrongCredentials  = NewUserError(nil, "wrong credentials", "", "US-6")
)

type UserError struct {
	Err              error  `json:"cause"`
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
