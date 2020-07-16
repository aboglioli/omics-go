package roles

import "omics/pkg/shared/errors"

var (
	Err         = errors.App.New().Path("security.domain.roles")
	ErrNotFound = Err.Code("not_found")
)
