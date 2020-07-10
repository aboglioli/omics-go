package security

import (
	"context"

	"omics/pkg/errors"
)

type TokenID string
type Token string

func NewTokenID() string {
	return "token-id"
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
