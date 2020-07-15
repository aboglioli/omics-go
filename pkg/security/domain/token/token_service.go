package token

import (
	"context"
	"fmt"

	"omics/pkg/shared/cache"
)

type TokenService struct {
	cache cache.Cache
	enc   TokenEncoder
}

func NewTokenService(cache cache.Cache, tokenEncoder TokenEncoder) *TokenService {
	return &TokenService{
		cache: cache,
		enc:   tokenEncoder,
	}
}

func (s *TokenService) Create(ctx context.Context, data Data) (Token, error) {
	tokenID := NewTokenID()

	token, err := s.enc.Encode(tokenID)
	if err != nil {
		return "", ErrToken.Code("create").Wrap(err)
	}

	if err := s.cache.Set(ctx, fmt.Sprintf("token:%s", tokenID), data); err != nil {
		return "", ErrToken.Code("create").Wrap(err)
	}

	return token, nil
}

func (s *TokenService) Validate(ctx context.Context, token Token) (Data, error) {
	tokenID, err := s.enc.Decode(token)
	if err != nil {
		return nil, ErrToken.Code("validate").Wrap(err)
	}

	rawData, err := s.cache.Get(ctx, fmt.Sprintf("token:%s", tokenID))
	if err != nil {
		return nil, ErrToken.Code("validate").Wrap(err)
	}

	if data, ok := rawData.(Data); ok {
		return data, nil
	}

	return nil, ErrToken.Code("validate")
}

func (s *TokenService) Invalidate(ctx context.Context, token Token) error {
	tokenID, err := s.enc.Decode(token)
	if err != nil {
		return ErrToken.Code("invalidate").Wrap(err)
	}

	if err := s.cache.Delete(ctx, fmt.Sprintf("token:%s", tokenID)); err != nil {
		return ErrToken.Code("invalidate").Wrap(err)
	}

	return nil
}
