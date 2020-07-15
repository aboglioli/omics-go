package users

type PasswordValidator struct{}

func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{}
}

func (pv *PasswordValidator) Validate(password string) error {
	return nil
}
