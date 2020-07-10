package token

import (
	"context"
	"fmt"

	"omics/pkg/common/cache"
	"omics/pkg/common/errors"
	"omics/pkg/security/domain/users"
)

type TokenService interface {
	Create(ctx context.Context, user *users.User) (Token, error)
	Validate(ctx context.Context, token Token) (*users.User, error)
	Update(ctx context.Context, token Token, user *users.User) error
	Invalidate(ctx context.Context, token Token) error
}

type tokenService struct {
	cache cache.Cache
	enc   TokenEncoder
}

func NewTokenService(cache cache.Cache, tokenEncoder TokenEncoder) TokenService {
	return &tokenService{
		cache: cache,
		enc:   tokenEncoder,
	}
}

func (s *tokenService) Create(ctx context.Context, user *users.User) (Token, error) {
	tokenID := NewTokenID()

	token, err := s.enc.Encode(tokenID)
	if err != nil {
		return "", errors.ErrTODO
	}

	if err := s.cache.Set(ctx, fmt.Sprintf("token:%s", tokenID), user); err != nil {
		return "", errors.ErrTODO
	}

	return token, nil
}

func (s *tokenService) Validate(ctx context.Context, token Token) (*users.User, error) {
	tokenID, err := s.enc.Decode(token)
	if err != nil {
		return nil, errors.ErrTODO
	}

	rawUser, err := s.cache.Get(ctx, fmt.Sprintf("token:%s", tokenID))
	if err != nil {
		return nil, errors.ErrTODO
	}

	if user, ok := rawUser.(*users.User); ok {
		return user, nil
	}

	return nil, errors.ErrTODO
}

func (s *tokenService) Update(ctx context.Context, token Token, user *users.User) error {
	tokenID, err := s.enc.Decode(token)
	if err != nil {
		return errors.ErrTODO
	}

	if _, err := s.cache.Get(ctx, fmt.Sprintf("token:%s", tokenID)); err != nil {
		return errors.ErrTODO
	}

	if err := s.cache.Set(ctx, fmt.Sprintf("token:%s", tokenID), user); err != nil {
		return errors.ErrTODO
	}

	return nil
}

func (s *tokenService) Invalidate(ctx context.Context, token Token) error {
	tokenID, err := s.enc.Decode(token)
	if err != nil {
		return errors.ErrTODO
	}

	if err := s.cache.Delete(ctx, fmt.Sprintf("token:%s", tokenID)); err != nil {
		return errors.ErrTODO
	}

	return nil
}
