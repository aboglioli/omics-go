package users

import "omics/pkg/common/errors"

var (
	ErrUsers = errors.Internal.New().Path("security.domain.users")
)
