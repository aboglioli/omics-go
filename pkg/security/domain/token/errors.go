package token

import "omics/pkg/shared/errors"

var (
	ErrToken = errors.Internal.New().Path("security.domain.token")
)
