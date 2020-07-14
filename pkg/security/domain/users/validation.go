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
	userID models.ID
	code   ValidationCode
}

func NewValidation(userID models.ID) *Validation {
	return &Validation{
		userID: userID,
		code:   NewValidationCode(),
	}
}

func (v *Validation) Validate(user *User, code string) error {
	if !user.ID().Equals(v.userID) {
		return ErrUsers
	}

	if !v.code.Equals(code) {
		return ErrUsers
	}

	user.validate()

	return nil
}
