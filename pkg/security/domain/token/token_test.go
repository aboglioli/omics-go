package token_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"omics/pkg/security/domain/token"
)

func TestTokenFromContext(t *testing.T) {
	tests := []struct {
		aCtx   context.Context
		rToken token.Token
		rErr   error
	}{{
		context.WithValue(context.Background(), "authToken", "#123#"),
		token.Token("#123#"),
		nil,
	}, {
		context.WithValue(context.Background(), "token", "#123#"),
		token.Token(""),
		token.ErrToken.Code("token_from_context"),
	}}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			rToken, rErr := token.TokenFromContext(test.aCtx)

			if !errors.Is(rErr, test.rErr) {
				t.Errorf("\nExp: %s\nAct: %s", test.rErr, rErr)
			} else if test.rToken != rToken {
				t.Errorf("\nExp: %s\nAct: %s", test.rToken, rToken)
			}
		})
	}
}
