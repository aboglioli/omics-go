package publications

type Name struct {
	name string
}

func NewName(name string) (Name, error) {
	if len(name) < 2 {
		return Name{}, ErrValidation.AddContext("name", "too_short")
	}

	if len(name) > 64 {
		return Name{}, ErrValidation.AddContext("name", "too_long")
	}

	return Name{name}, nil
}
