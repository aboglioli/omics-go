package security

import "omics/pkg/errors"

var (
	ErrNull = errors.New(errors.APPLICATION, "null")
)
