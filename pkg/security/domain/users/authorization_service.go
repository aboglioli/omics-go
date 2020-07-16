package users

import (
	"context"

	"omics/pkg/security/domain/roles"
)

type Guard func(*User) bool

func DefaultGuard() bool {
	return true
}

type AuthorizationService struct {
	userRepo UserRepository
	roleRepo roles.RoleRepository
}

func (s *AuthorizationService) UserHasPermissions(
	ctx context.Context,
	user *User,
	permissions,
	module string,
	guard Guard,
) bool {
	role, err := s.roleRepo.FindByCode(ctx, user.RoleCode())
	if role == nil || err != nil {
		return false
	}

	if role.Is(roles.ADMIN) {
		return true
	}

	return role.HasPermissions(permissions, module) && guard(user)
}
