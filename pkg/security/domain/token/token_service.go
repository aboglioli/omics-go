//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks
package token

import (
	"context"
	"fmt"

	"omics/pkg/common/cache"
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
		return "", ErrToken.Code("create").Wrap(err)
	}

	if err := s.cache.Set(ctx, fmt.Sprintf("token:%s", tokenID), user); err != nil {
		return "", ErrToken.Code("create").Wrap(err)
	}

	return token, nil
}

func (s *tokenService) Validate(ctx context.Context, token Token) (*users.User, error) {
	tokenID, err := s.enc.Decode(token)
	if err != nil {
		return nil, ErrToken.Code("validate").Wrap(err)
	}

	rawUser, err := s.cache.Get(ctx, fmt.Sprintf("token:%s", tokenID))
	if err != nil {
		return nil, ErrToken.Code("validate").Wrap(err)
	}

	if user, ok := rawUser.(*users.User); ok {
		return user, nil
	}

	return nil, ErrToken.Code("validate")
}

func (s *tokenService) Update(ctx context.Context, token Token, user *users.User) error {
	tokenID, err := s.enc.Decode(token)
	if err != nil {
		return ErrToken.Code("update").Wrap(err)
	}

	if _, err := s.cache.Get(ctx, fmt.Sprintf("token:%s", tokenID)); err != nil {
		return ErrToken.Code("update").Wrap(err)
	}

	if err := s.cache.Set(ctx, fmt.Sprintf("token:%s", tokenID), user); err != nil {
		return ErrToken.Code("update").Wrap(err)
	}

	return nil
}

func (s *tokenService) Invalidate(ctx context.Context, token Token) error {
	tokenID, err := s.enc.Decode(token)
	if err != nil {
		return ErrToken.Code("invalidate").Wrap(err)
	}

	if err := s.cache.Delete(ctx, fmt.Sprintf("token:%s", tokenID)); err != nil {
		return ErrToken.Code("invalidate").Wrap(err)
	}

	return nil
}
