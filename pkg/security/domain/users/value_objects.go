package users

import "omics/pkg/shared/models"

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

type Username struct {
	Username string `validate:"required,min=4,max=24"`
}

func NewUsername(username string) (Username, error) {
	u := Username{username}
	if err := models.Validate(&u); err != nil {
		return Username{}, ErrValidation.Merge(err)
	}
	return u, nil
}

type Email struct {
	Email string `validate:"required,min=5,email"`
}

func NewEmail(email string) (Email, error) {
	e := Email{email}
	if err := models.Validate(&e); err != nil {
		return Email{}, nil
	}
	return e, nil
}
