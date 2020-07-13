package users

import "omics/pkg/shared/errors"

var (
	ErrUsers = errors.Internal.New().Path("security.domain.users")
)
