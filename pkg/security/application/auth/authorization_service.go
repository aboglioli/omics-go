package auth

import (
	"context"

	"omics/pkg/security/domain/token"
	"omics/pkg/security/domain/users"
	"omics/pkg/shared/models"
)

type AuthorizationService interface {
	GetUserByID(ctx context.Context, userID models.ID) (*users.User, error)
	GetUserByToken(ctx context.Context, t token.Token) (*users.User, error)
	UserHasRole(ctx context.Context, userID models.ID, role string) error
	UserHasPermissions(ctx context.Context, userID models.ID, permissions, module string) error
}

type authorizationService struct {
	tokenServ token.TokenService
	userRepo  users.UserRepository
}

func NewAuthorizationService(tokenServ token.TokenService, userRepo users.UserRepository) AuthorizationService {
	return &authorizationService{
		tokenServ: tokenServ,
		userRepo:  userRepo,
	}
}

func (s *authorizationService) GetUserByID(ctx context.Context, userID models.ID) (*users.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUnauthorized.Wrap(err)
	}

	return user, nil
}

func (s *authorizationService) GetUserByToken(ctx context.Context, t token.Token) (*users.User, error) {
	user, err := s.tokenServ.Validate(ctx, t)
	if err != nil {
		return nil, ErrUnauthorized.Wrap(err)
	}

	return user, nil
}

func (s *authorizationService) UserHasRole(ctx context.Context, userID models.ID, role string) error {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.HasRole(role) {
		return nil
	}

	return ErrUnauthorized
}

func (s *authorizationService) UserHasPermissions(
	ctx context.Context,
	userID models.ID,
	permissions string,
	module string,
) error {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.HasPermissions(permissions, module) {
		return nil
	}

	return ErrUnauthorized
}
