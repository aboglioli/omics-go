package users

import "context"

type AuthenticationService interface {
	Authenticate(ctx context.Context, usernameOrEmail, password string) (*User, error)
}

type authenticationService struct {
	userRepo UserRepository
	userServ UserService
}

func (s *authenticationService) Authenticate(ctx context.Context, usernameOrEmail, password string) (*User, error) {
	user, err := s.userRepo.FindByUsername(ctx, usernameOrEmail)
	if user == nil || err != nil {
		user, err = s.userRepo.FindByEmail(ctx, usernameOrEmail)
		if user == nil || err != nil {
			return nil, ErrUnauthorized
		}
	}

	if !s.userServ.ComparePassword(user, password) {
		return nil, ErrUnauthorized
	}

	user.wasAuthenticated()

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, ErrUnauthorized
	}

	return user, nil
}
