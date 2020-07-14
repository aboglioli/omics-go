//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks
package users

import (
	"context"
	"omics/pkg/shared/models"
)

type UserRepository interface {
	NextID() models.ID
	FindByID(ctx context.Context, userID models.ID) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Save(ctx context.Context, user *User) error
}
