package token

import (
	"context"

	"omics/pkg/common/errors"

	"github.com/google/uuid"
)

type TokenID string
type Token string

func NewTokenID() TokenID {
	uuid := uuid.New().String()
	return TokenID(uuid)
}

func TokenFromContext(ctx context.Context) (Token, error) {
	if tokenStr, ok := ctx.Value("authToken").(string); ok {
		return Token(tokenStr), nil
	}
	return "", errors.ErrTODO
}

func (t Token) String() string {
	return string(t)
}
