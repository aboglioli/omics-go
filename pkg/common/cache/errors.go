package cache

import "omics/pkg/common/errors"

var (
	ErrCache = errors.Internal.New().Path("common.cache")
)
