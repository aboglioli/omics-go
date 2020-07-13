package roles

import "omics/pkg/shared/errors"

var (
	ErrRoles = errors.Internal.New().Path("security.domain.roles")
)
