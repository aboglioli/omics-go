package token

import (
	"context"

	"github.com/google/uuid"
)

type TokenID string
type Token string

func NewTokenID() TokenID {
	uuid := uuid.New().String()
	return TokenID(uuid)
}

func FromContext(ctx context.Context) (Token, error) {
	if tokenStr, ok := ctx.Value("authToken").(string); ok {
		return Token(tokenStr), nil
	}
	return "", ErrToken.Code("token_from_context")
}

func (t Token) String() string {
	return string(t)
}
