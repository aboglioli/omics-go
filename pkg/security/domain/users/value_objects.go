package users

import "omics/pkg/shared/models"

type Identity struct {
	Username string `validate:"required,min=4,max=24"`
	Email    string `validate:"required,min=5,email"`
}

func NewIdentity(username, email string) (Identity, error) {
	id := Identity{username, email}
	if err := models.Validate(&id); err != nil {
		return Identity{}, ErrValidation.Merge(err)
	}

	return id, nil
}

type Fullname struct {
	Name     string `validate:"required,min=2,max=16"`
	Lastname string `validate:"required,min=2,max=16"`
}

func NewFullname(name, lastname string) (Fullname, error) {
	n := Fullname{name, lastname}
	if err := models.Validate(&n); err != nil {
		return Fullname{}, ErrValidation.Merge(err)
	}
	return n, nil
}
