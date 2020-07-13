package cache

import "omics/pkg/shared/errors"

var (
	ErrCache = errors.Internal.New().Path("shared.cache")
)
