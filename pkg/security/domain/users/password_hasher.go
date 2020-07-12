//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks
package users

type PasswordHasher interface {
	Hash(plainPassword string) (string, error)
	Compare(hashedPassword, plainPassword string) bool
}
