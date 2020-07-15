package publications

import "omics/pkg/shared/errors"

var (
	Err           = errors.App.New().Path("publishing.domain.publications")
	ErrValidation = Err.Code("validation")
)
