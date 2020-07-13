package auth

import "context"

type AuthenticationService interface {
	Authenticate(ctx context.Context, usernameOrEmail, password string)
}
