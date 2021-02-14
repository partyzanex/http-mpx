package api

// Error represents a custom error for HTTP-response.
type Error struct {
	Code     int    `json:"-"`
	Internal error  `json:"-"`
	Message  string `json:"message"`
}

// Error implements the error interface.
func (e Error) Error() string {
	return e.Message
}

// SetInternal sets Internal field.
func (e *Error) SetInternal(err error) *Error {
	e.Internal = err
	return e
}

// NewError creates *Error with http status code and message.
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
