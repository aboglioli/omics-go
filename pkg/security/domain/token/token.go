package token

import (
	"context"

	"github.com/google/uuid"
)

type TokenID string

func NewTokenID() TokenID {
	uuid := uuid.New().String()
	return TokenID(uuid)
}

type Token string

func (t Token) String() string {
	return string(t)
}

type Data map[string]string

func NewData(userID string) Data {
	return Data{
		"user_id": userID,
	}
}

func FromContext(ctx context.Context) (Token, error) {
	if tokenStr, ok := ctx.Value("authToken").(string); ok {
		return Token(tokenStr), nil
	}
	return "", ErrToken.Code("token_from_context")
}
