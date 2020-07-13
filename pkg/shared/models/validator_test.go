package models

import (
	"reflect"
	"strconv"
	"testing"

	"omics/pkg/shared/errors"
)

func TestValidator(t *testing.T) {
	type data struct {
		Username string  `validate:"required,min=4,max=32"`
		Password string  `validate:"required"`
		Email    string  `validate:"required,min=5,max=64,email"`
		Name     string  `validate:"min=2,max=32,alphaspaces"`
		Lastname *string `validate:"omitempty,min=4,max=6"`
	}

	tests := []struct {
		aData *data
		rErr  error
	}{{
		&data{},
		ErrValidation.Context(errors.Context{
			"username": "required",
			"password": "required",
			"email":    "required",
			"name":     "min",
		}),
	}, {
		&data{
			Username: "adm",
			Password: "123",
			Email:    "asd.com",
			Name:     "n",
			Lastname: NewString("las"),
		},
		ErrValidation.Context(errors.Context{
			"username": "min",
			"email":    "email",
			"name":     "min",
			"lastname": "min",
		}),
	}, {
		&data{
			Username: "admin",
			Password: "123",
			Email:    "asd@asd.com",
			Name:     "name",
			Lastname: NewString("lastname"),
		},
		ErrValidation.Context(errors.Context{
			"lastname": "max",
		}),
	}, {
		&data{
			Username: "admin",
			Password: "123",
			Email:    "asd@asd.com",
			Name:     "name",
			Lastname: NewString("last"),
		},
		nil,
	}}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			rErr := Validate(test.aData)

			if !reflect.DeepEqual(test.rErr, rErr) {
				t.Errorf("\nExp: %s\nAct: %s", test.rErr, rErr)
			}
		})
	}
}
