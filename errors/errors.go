package errors

import (
	"fmt"
	"github.com/nit-app/nit-backend/models/status"
)

type Error struct {
	typ string
	err error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.typ, e.err)
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Type() string {
	return e.typ
}

func (e *Error) MakeResponse() (statusCode, reason string) {
	// internal server error is the special case of error handling where we don't report reason to user
	if !status.Codes[e.typ].ExposeText || e.err == nil {
		return e.typ, ""
	}

	return e.typ, e.err.Error()
}

func New(typ string, err error) *Error {
	return &Error{typ: typ, err: err}
}
