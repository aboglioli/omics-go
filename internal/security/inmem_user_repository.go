package security

import (
	"context"
	"omics/pkg/models"
)

type inmemUserRepository struct {
	users []*User
}

func (r *inmemUserRepository) FindByID(ctx context.Context, userID models.ID) (*User, error) {
	for _, user := range r.users {
		if user.ID == userID {
			return user, nil
		}
	}
	return nil, ErrNull
}

func (r *inmemUserRepository) FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*User, error) {
	for _, user := range r.users {
		if user.Username == usernameOrEmail || user.Email == usernameOrEmail {
			return user, nil
		}
	}
	return nil, ErrNull
}

func (r *inmemUserRepository) Save(ctx context.Context, user *User) error {
	if user.ID.Str() == "" {
		r.users = append(r.users, user)
		return nil
	}

	for i, u := range r.users {
		if u.ID == user.ID {
			r.users[i] = user
			return nil
		}
	}

	return ErrNull
}
