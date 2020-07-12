package events

import "omics/pkg/common/errors"

var (
	ErrEvents = errors.Internal.New().Path("common.events")
)
