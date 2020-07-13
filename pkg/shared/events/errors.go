package events

import "omics/pkg/shared/errors"

var (
	ErrEvents = errors.Internal.New().Path("shared.events")
)
