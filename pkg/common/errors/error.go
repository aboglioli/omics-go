package errors

import (
	"encoding/json"
	"fmt"
)

type ErrorType string

const (
	INTERNAL    ErrorType = "internal"
	APPLICATION ErrorType = "application"
)

type Context map[string]interface{}

type Error struct {
	kind    ErrorType
	code    string
	path    string
	status  int
	message string
	context Context
	cause   error
}

func New(t ErrorType, c string) Error {
	return Error{
		kind:    t,
		code:    c,
		context: make(Context),
	}
}

func (e Error) Path(path string) Error {
	e.path = path
	return e
}

func (e Error) Status(status int) Error {
	e.status = status
	return e
}

func (e Error) Message(msg string, args ...interface{}) Error {
	e.message = fmt.Sprintf(msg, args...)
	return e
}

func (e Error) Context(ctx Context) Error {
	newCtx := make(Context)
	for k, v := range e.context {
		newCtx[k] = v
	}
	for k, v := range ctx {
		newCtx[k] = v
	}
	e.context = newCtx
	return e
}

func (e Error) AddContext(k string, v interface{}) Error {
	ctx := make(Context)
	for k, v := range e.context {
		ctx[k] = v
	}
	ctx[k] = v
	e.context = ctx
	return e
}

func (e Error) Cause(cause error) Error {
	e.cause = cause
	return e
}

// DisplayError
type displayError struct {
	Kind    ErrorType `json:"type,omitempty"`
	Code    string    `json:"code,omitempty"`
	Path    string    `json:"path,omitempty"`
	Status  int       `json:"status,omitempty"`
	Message string    `json:"message,omitempty"`
	Context Context   `json:"context,omitempty"`
	Cause   string    `json:"cause,omitempty"`
}

func (e Error) display() *displayError {
	err := &displayError{
		Kind:    e.kind,
		Code:    e.code,
		Path:    e.path,
		Status:  e.status,
		Message: e.message,
		Context: e.context,
	}
	if e.cause != nil {
		err.Cause = e.cause.Error()
	}
	return err
}

func (e Error) Error() string {
	str, err := json.Marshal(e.display())
	if err != nil {
		return ""
	}
	return string(str)
}

func (e Error) String() string {
	return e.Error()
}
