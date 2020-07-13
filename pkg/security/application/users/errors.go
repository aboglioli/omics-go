package users

import "omics/pkg/shared/errors"

var (
	ErrUsers        = errors.App.New().Path("security.application.users")
	ErrNotFound     = ErrUsers.Code("not_found")
	ErrUnauthorized = ErrUsers.Code("unauthorized")
)
