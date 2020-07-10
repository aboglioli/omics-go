package security

import (
	"context"

	"omics/pkg/errors"
	"omics/pkg/models"
)

type inmemUserRepository struct {
	users map[models.ID]*User
}

func NewInMemUserRepository() *inmemUserRepository {
	return &inmemUserRepository{
		users: make(map[models.ID]*User),
	}
}

func (r *inmemUserRepository) FindByID(ctx context.Context, userID models.ID) (*User, error) {
	if user, ok := r.users[userID]; ok {
		return user, nil
	}
	return nil, errors.ErrTODO
}

func (r *inmemUserRepository) FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*User, error) {
	for _, user := range r.users {
		if user.Username == usernameOrEmail || user.Email == usernameOrEmail {
			return user, nil
		}
	}
	return nil, errors.ErrTODO
}

func (r *inmemUserRepository) Save(ctx context.Context, user *User) error {
	r.users[user.ID] = user
	return nil
}
