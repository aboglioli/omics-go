package users

import "omics/pkg/shared/errors"

var (
	ErrUsers      = errors.App.New().Path("security.domain.users")
	ErrValidation = ErrUsers.Code("validation")
)
