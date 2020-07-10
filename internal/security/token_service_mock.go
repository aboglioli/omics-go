package security

import (
	"context"
	"fmt"

	"omics/pkg/errors"
)

type fakeTokenService struct {
	Users map[Token]*User
}

func FakeTokenService() *fakeTokenService {
	return &fakeTokenService{
		Users: make(map[Token]*User),
	}
}

func (s *fakeTokenService) Create(ctx context.Context, user *User) (Token, error) {
	token := Token(fmt.Sprintf("token:%d", user.ID))
	s.Users[token] = user
	return token, nil
}

func (s *fakeTokenService) Validate(ctx context.Context, token Token) (*User, error) {
	if user, ok := s.Users[token]; ok {
		return user, nil
	}
	return nil, errors.ErrTODO
}

func (s *fakeTokenService) Invalidate(ctx context.Context, token Token) error {
	if _, ok := s.Users[token]; !ok {
		return errors.ErrTODO
	}
	delete(s.Users, token)
	return nil
}
