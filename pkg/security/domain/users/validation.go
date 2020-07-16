package users

import (
	"omics/pkg/shared/models"

	"github.com/google/uuid"
)

type ValidationCode struct {
	uniqueCode string
}

func NewValidationCode() ValidationCode {
	uuid := uuid.New().String()
	return ValidationCode{uuid}
}

func (vc ValidationCode) Equals(code string) bool {
	return vc.uniqueCode == code
}

type Validation struct {
	id     models.ID
	userID models.ID
	code   ValidationCode
	active bool
}

func NewValidation(userID models.ID) *Validation {
	return &Validation{
		userID: userID,
		code:   NewValidationCode(),
	}
}

func (v *Validation) Validate(user *User, code string) error {
	if !user.Equals(v.userID) {
		return ErrValidation.AddContext("user", "not_equal")
	}

	if !v.code.Equals(code) {
		return ErrValidation.AddContext("code", "invalid")
	}

	user.validate()

	return nil
}
