package users

type PasswordValidator interface {
	Validate(password string) error
}
