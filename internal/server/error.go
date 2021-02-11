package server

type Error struct {
	Code     int    `json:"-"`
	Internal error  `json:"-"`
	Message  string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func (e *Error) SetInternal(err error) *Error {
	e.Internal = err
	return e
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
