package token

import "omics/pkg/common/errors"

var (
	ErrToken = errors.Internal.New().Path("security.domain.token")
)
