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

type AuthenticationService interface {
	Authenticate(ctx context.Context, cmd *AuthenticateCommand) (token.Token, error)
	Deauthenticate(ctx context.Context, t token.Token) error
}

type authenticationService struct {
	userRepo       users.UserRepository
	tokenServ      token.TokenService
	passwordHasher users.PasswordHasher
}

func (s *authenticationService) Authenticate(
	ctx context.Context,
	cmd *AuthenticateCommand,
) (token.Token, error) {
	user, err := s.userRepo.FindByUsername(ctx, cmd.UsernameOrEmail)
	if err != nil {
		user, err = s.userRepo.FindByEmail(ctx, cmd.UsernameOrEmail)
		if err != nil {
			return token.Token(""), ErrAuthentication.AddContext("username", "invalid").Wrap(err)
		}
	}

	if !user.ComparePassword(cmd.Password, s.passwordHasher) {
		return token.Token(""), ErrAuthentication.AddContext("password", "invalid")
	}

	t, err := s.tokenServ.Create(ctx, user)
	if err != nil {
		return token.Token(""), ErrAuthentication.Wrap(err)
	}

	return t, nil
}

func (s *authenticationService) Deauthenticate(ctx context.Context, t token.Token) error {
	if err := s.tokenServ.Invalidate(ctx, t); err != nil {
		return ErrUnauthorized.Wrap(err)
	}

	return nil
}
