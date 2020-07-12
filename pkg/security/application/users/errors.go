package users

import "omics/pkg/common/errors"

var (
	ErrUsers        = errors.App.New().Path("security.application.users")
	ErrNotFound     = ErrUsers.Code("not_found")
	ErrUnauthorized = ErrUsers.Code("unauthorized")
)
