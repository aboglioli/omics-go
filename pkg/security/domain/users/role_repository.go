//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks
package users

import "context"

type RoleRepository interface {
	FindByCode(ctx context.Context, code string) (Role, error)
}
