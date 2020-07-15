package users

import (
	"context"

	"omics/pkg/security/domain/roles"
)

type AuthorizationService struct {
	userRepo UserRepository
	roleRepo roles.RoleRepository
}

func (s *AuthenticationService) UserHasPermissions(ctx context.Context, user *User, permissions, module string) bool {
	role, err := s.roleRepo.FindByCode(ctx, user.RoleCode())
	if role == nil || err != nil {
		return false
	}

	return role.HasPermissions(permissions, module)
}
