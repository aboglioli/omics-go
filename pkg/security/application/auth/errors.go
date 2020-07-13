package auth

import "omics/pkg/shared/errors"

var (
	ErrUnauthorized   = errors.App.New().Path("security.application.auth").Code("unauthorized")
	ErrAuthentication = errors.App.New().Path("security.application.auth").Code("authentication")
)
