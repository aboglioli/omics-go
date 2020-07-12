package roles

import "omics/pkg/common/errors"

var (
	ErrRoles = errors.Internal.New().Path("security.domain.roles")
)
