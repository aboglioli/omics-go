package models

import (
	"regexp"

	"omics/pkg/common/errors"

	govalidator "github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

var (
	ErrValidation = errors.App.New().Code("validation")
)

var (
	alphaWithSpacesRE  = regexp.MustCompile("^[a-zA-Záéíóú ]*$")
	alphaNumWithDashRE = regexp.MustCompile("^[a-zA-Z0-9-]*$")
)

func alphaWithSpaces(fl govalidator.FieldLevel) bool {
	str := fl.Field().String()
	if str == "invalid" {
		return false
	}

	return alphaWithSpacesRE.MatchString(str)
}

func alphaNumWithDash(fl govalidator.FieldLevel) bool {
	str := fl.Field().String()
	if str == "invalid" {
		return false
	}

	return alphaNumWithDashRE.MatchString(str)
}

func Validate(s interface{}) error {
	validator := govalidator.New()
	validator.RegisterValidation("alphaspaces", alphaWithSpaces)
	validator.RegisterValidation("alphanumdash", alphaNumWithDash)

	if err := validator.Struct(s); err != nil {
		fieldErrors := ErrValidation

		if errs, ok := err.(govalidator.ValidationErrors); ok {
			for _, err := range errs {
				field := strcase.ToLowerCamel(err.Field())
				fieldErrors = fieldErrors.AddContext(field, err.Tag())
			}
			return fieldErrors
		}

		return fieldErrors
	}
	return nil
}
