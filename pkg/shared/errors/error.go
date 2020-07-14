package errors

import (
	"encoding/json"
	"fmt"
)

type ErrorKind string

const (
	Raw      ErrorKind = "raw"
	Internal ErrorKind = "internal"
	App      ErrorKind = "application"
	TODO     ErrorKind = "todo"
	Test     ErrorKind = "test"
)

func (et ErrorKind) New() Error {
	return Error{
		kind:    et,
		context: make(Context),
	}
}

type Context map[string]interface{}

type Error struct {
	kind    ErrorKind
	code    string
	path    string
	status  int
	message string
	context Context
	cause   error
}

func (e Error) Code(code string) Error {
	e.code = code
	return e
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

func (e Error) Merge(err error) Error {
	if err, ok := err.(Error); ok {
		return e.Context(err.context)
	}
	return e.Message(err.Error())
}

func (e Error) ContextLen() int {
	return len(e.context)
}

// error interface implementations
func (e Error) Wrap(cause error) Error {
	e.cause = cause
	return e
}

func (e Error) Unwrap() error {
	return e.cause
}

func (e Error) Is(err error) bool {
	if err, ok := err.(Error); ok {
		return e.code == err.code
	}
	return false
}

func (e Error) Error() string {
	str, err := json.Marshal(Display(e, true))
	if err != nil {
		return ""
	}
	return string(str)
}

func (e Error) String() string {
	return e.Error()
}
