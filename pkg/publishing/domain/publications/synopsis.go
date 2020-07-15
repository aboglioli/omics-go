package publications

type Synopsis struct {
	synopsis string
}

func NewSynopsis(synopsis string) (Synopsis, error) {
	if len(synopsis) < 21 {
		return Synopsis{}, ErrValidation.AddContext("synopsis", "too_short")
	}

	if len(synopsis) > 512 {
		return Synopsis{}, ErrValidation.AddContext("synopsis", "too_long")
	}

	return Synopsis{synopsis}, nil
}
