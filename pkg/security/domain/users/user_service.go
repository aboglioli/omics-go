package users

import "context"

type UserService interface {
	Available(ctx context.Context, username, email string) error
}

type userService struct {
	userRepo UserRepository
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
