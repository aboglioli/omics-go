package security

import "context"

type TokenService interface {
	Create(ctx context.Context, user *User) (Token, error)
	Validate(ctx context.Context, token Token) (*User, error)
	Invalidate(ctx context.Context, token Token) error
}
