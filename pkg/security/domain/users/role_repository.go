package users

import "context"

// TODO: create own entity for role
type RoleRepository interface {
	FindByCode(ctx context.Context, code string) (Role, error)
}
