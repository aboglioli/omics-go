package persistence

import (
	"context"

	"omics/pkg/common/errors"
	"omics/pkg/common/models"
	"omics/pkg/security/domain/users"
)

type inmemUserRepository struct {
	users map[models.ID]*users.User
}

func NewInMemUserRepository() *inmemUserRepository {
	return &inmemUserRepository{
		users: make(map[models.ID]*users.User),
	}
}

func (r *inmemUserRepository) FindByID(ctx context.Context, userID models.ID) (*users.User, error) {
	if user, ok := r.users[userID]; ok {
		return user, nil
	}
	return nil, errors.ErrTODO
}

func (r *inmemUserRepository) FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*users.User, error) {
	for _, user := range r.users {
		if user.Username == usernameOrEmail || user.Email == usernameOrEmail {
			return user, nil
		}
	}
	return nil, errors.ErrTODO
}

func (r *inmemUserRepository) Save(ctx context.Context, user *users.User) error {
	r.users[user.ID] = user
	return nil
}
