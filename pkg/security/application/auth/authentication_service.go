package auth

import (
	"context"

	"omics/pkg/security/domain/token"
	"omics/pkg/security/domain/users"
)

type AuthenticateCommand struct {
	UsernameOrEmail string `json:"username"`
	Password        string `json:"password"`
}

type AuthenticateResponse struct {
	AuthToken token.Token `json:"auth_token"`
}

type AuthenticationService interface {
	Authenticate(ctx context.Context, cmd *AuthenticateCommand) (*AuthenticateResponse, error)
	Deauthenticate(ctx context.Context) error
}

type authenticationService struct {
	userRepo       users.UserRepository
	tokenServ      token.TokenService
	passwordHasher users.PasswordHasher
}

func (s *authenticationService) Authenticate(
	ctx context.Context,
	cmd *AuthenticateCommand,
) (*AuthenticateResponse, error) {
	user, err := s.userRepo.FindByUsernameOrEmail(ctx, cmd.UsernameOrEmail)
	if err != nil {
		return nil, ErrAuthentication.Wrap(err)
	}

	if !user.ComparePassword(cmd.Password, s.passwordHasher) {
		return nil, ErrAuthentication.AddContext("password", "mismatch")
	}

	t, err := s.tokenServ.Create(ctx, user)
	if err != nil {
		return nil, ErrAuthentication.Wrap(err)
	}

	return &AuthenticateResponse{
		AuthToken: t,
	}, nil
}

func (s *authenticationService) Deauthenticate(ctx context.Context) error {
	t, err := token.TokenFromContext(ctx)
	if err != nil {
		return ErrUnauthorized.Wrap(err)
	}

	if err := s.tokenServ.Invalidate(ctx, t); err != nil {
		return ErrUnauthorized.Wrap(err)
	}

	return nil
}
