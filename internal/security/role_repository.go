package security

import "context"

type RoleRepository interface {
	FindByCode(ctx context.Context, code string) (Role, error)
}
