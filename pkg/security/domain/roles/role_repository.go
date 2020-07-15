//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks
package roles

import "context"

type RoleRepository interface {
	FindByCode(ctx context.Context, code string) (*Role, error)
}
