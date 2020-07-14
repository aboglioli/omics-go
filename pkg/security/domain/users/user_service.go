package users

import "context"

type UserService interface {
	Available(ctx context.Context, username, email string) error
	ChangePassword(user *User, oldPassword, newPassword string) error
	ComparePassword(user *User, password string) bool
}

type userService struct {
	userRepo          UserRepository
	userServ          UserService
	passwordHasher    PasswordHasher
	passwordValidator PasswordValidator
}

func (s *userService) Available(ctx context.Context, username, email string) error {
	errs := ErrUsers
	if user, err := s.userRepo.FindByUsername(ctx, username); user != nil || err == nil {
		errs = errs.AddContext("username", "not_available")
	}

	if user, err := s.userRepo.FindByEmail(ctx, email); user != nil || err == nil {
		errs = errs.AddContext("email", "not_available")
	}

	if errs.ContextLen() > 0 {
		return errs
	}

	return nil
}

func (s *userService) ChangePassword(user *User, oldPassword, newPassword string) error {
	if user.password != "" && !s.passwordHasher.Compare(user.password, oldPassword) {
		return ErrUnauthorized
	}

	if err := s.passwordValidator.Validate(newPassword); err != nil {
		return err
	}

	hashedPassword, err := s.passwordHasher.Hash(newPassword)
	if err != nil {
		return ErrUsers.Wrap(err)
	}

	user.password = hashedPassword

	return nil
}

func (s *userService) ComparePassword(user *User, password string) bool {
	return s.passwordHasher.Compare(user.password, password)
}
