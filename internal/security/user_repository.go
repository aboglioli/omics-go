package security

import (
	"context"
	"omics/pkg/models"
)

type UserRepository interface {
	FindByID(ctx context.Context, userID models.ID) (*User, error)
	FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*User, error)
	Save(ctx context.Context, user *User) error
}
