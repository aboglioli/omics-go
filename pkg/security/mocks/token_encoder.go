package mocks

import (
	"fmt"
	"strings"

	"omics/pkg/common/errors"
	"omics/pkg/security/domain/token"
)

type fakeTokenEncoder struct{}

func FakeTokenEncoder() *fakeTokenEncoder {
	return &fakeTokenEncoder{}
}

func (e *fakeTokenEncoder) Encode(tokenID token.TokenID) (token.Token, error) {
	t := fmt.Sprintf("token:%s", tokenID)
	return token.Token(t), nil
}

func (e *fakeTokenEncoder) Decode(t token.Token) (token.TokenID, error) {
	tokenStr := t.String()
	if !strings.HasPrefix(tokenStr, "token:") {
		return "", errors.ErrTODO
	}

	tokenID := strings.TrimPrefix(tokenStr, "token:")
	return token.TokenID(tokenID), nil
}
