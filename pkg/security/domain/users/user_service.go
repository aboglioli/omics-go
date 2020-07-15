package users

import "context"

type UserService struct {
	passwordHasher    PasswordHasher
	passwordValidator *PasswordValidator
	userRepo          UserRepository
}

func NewUserService(
	passwordHasher PasswordHasher,
	passwordValidator *PasswordValidator,
	userRepo UserRepository,
) *UserService {
	return &UserService{
		passwordHasher:    passwordHasher,
		passwordValidator: passwordValidator,
	}
}

func (s *UserService) Available(ctx context.Context, username, email string) error {
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

func (s *UserService) ChangePassword(user *User, oldPassword, newPassword string) error {
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

func (s *UserService) ComparePassword(user *User, password string) bool {
	return s.passwordHasher.Compare(user.password, password)
}
