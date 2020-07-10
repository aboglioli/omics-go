package security

import (
	"fmt"
	"omics/pkg/errors"
	"strings"
)

type fakeTokenEncoder struct{}

func FakeTokenEncoder() *fakeTokenEncoder {
	return &fakeTokenEncoder{}
}

func (e *fakeTokenEncoder) Encode(tokenID TokenID) (Token, error) {
	token := fmt.Sprintf("token:%s", tokenID)
	return Token(token), nil
}

func (e *fakeTokenEncoder) Decode(token Token) (TokenID, error) {
	tokenStr := token.String()
	if !strings.HasPrefix(tokenStr, "token:") {
		return "", errors.ErrTODO
	}

	tokenID := strings.TrimPrefix(tokenStr, "token:")
	return TokenID(tokenID), nil
}
