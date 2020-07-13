//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks
package users

import (
	"context"
	"omics/pkg/common/models"
)

type UserRepository interface {
	NextID() models.ID
	FindByID(ctx context.Context, userID models.ID) (*User, error)
	FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*User, error)
	Save(ctx context.Context, user *User) error
}
