package token

import "omics/pkg/shared/errors"

var (
	Err = errors.Internal.New().Path("security.domain.token")
)
