package errors

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(&Error{})
}

func New(format string, args ...interface{}) *Error {
	return N(format, args...)
}

// create a new error message
func N(format string, args ...interface{}) *Error {
	return &Error{fmt.Sprintf(format, args...)}
}

// create a new error by wrapping an existing error
func W(err error) *Error {
	return N("underlying %s", err.Error())
}

type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}
