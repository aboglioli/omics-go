package users

import "omics/pkg/shared/errors"

var (
	Err             = errors.App.New().Path("security.domain.users")
	ErrValidation   = Err.Code("validation")
	ErrUnauthorized = Err.Code("unauthorized").Status(400)
	ErrNotFound     = Err.Code("not_found").Status(404)
)
