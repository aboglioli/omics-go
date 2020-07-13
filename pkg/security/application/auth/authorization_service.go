package auth

import (
	"context"
	"omics/pkg/security/domain/users"
	"omics/pkg/shared/models"
)

type AuthorizationService interface {
	GetUserFromCtx(ctx context.Context) (*users.User, error)

	UserHasRole(ctx context.Context, userID models.ID, role string) error
	UserHasRoleFromCtx(ctx context.Context, role string) error

	UserHasPermissions(ctx context.Context, userID models.ID, permission, module string) error
	UserHasPermissionsFromCtx(ctx context.Context, permission, module string) error
}
