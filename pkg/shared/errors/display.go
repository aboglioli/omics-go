package errors

import "encoding/json"

func Display(e error, includeInternal bool) *DisplayError {
	if e, ok := e.(Error); ok {
		if e.kind == Internal && !includeInternal {
			return nil
		}

		err := &DisplayError{
			Kind:    e.kind,
			Code:    e.code,
			Path:    e.path,
			Status:  e.status,
			Message: e.message,
		}

		if len(e.context) > 0 {
			err.Context = e.context
		}

		if e.cause != nil {
			if causeErr := Display(e.cause, includeInternal); causeErr != nil {
				err.Cause = causeErr
			}
		}

		return err
	}

	if !includeInternal {
		return nil
	}

	return &DisplayError{
		Kind:    Raw,
		Message: e.Error(),
	}
}

type DisplayError struct {
	Kind    ErrorKind `json:"kind,omitempty"`
	Code    string    `json:"code,omitempty"`
	Path    string    `json:"path,omitempty"`
	Status  int       `json:"status,omitempty"`
	Message string    `json:"message,omitempty"`
	Context Context   `json:"context,omitempty"`
	Cause   error     `json:"cause,omitempty"`
}

func (e *DisplayError) Error() string {
	str, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(str)
}

func (e *DisplayError) String() string {
	return e.Error()
}
