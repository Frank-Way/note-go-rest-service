package nerror

import "encoding/json"

var (
	ErrorNotFound   = NewNoteError(nil, "note not found", "", "NS-1")
	ErrorDuplicate  = NewNoteError(nil, "note already exists", "", "NS-2")
	ErrorRepository = NewNoteError(nil, "repository error", "", "NS-3")
	ErrorNoAuth     = NewNoteError(nil, "no authorized", "", "NS-4")
)

type NoteError struct {
	Err              error  `json:"-"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developer_message"`
	Code             string `json:"code"`
}

func (e *NoteError) Error() string {
	return e.Message
}

func (e *NoteError) Unwrap() error { return e.Err }

func (e *NoteError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewNoteError(err error, message, developerMessage, code string) *NoteError {
	return &NoteError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}

func systemError(err error) *NoteError {
	return NewNoteError(err, "internal system error", err.Error(), "NS-0")
}
