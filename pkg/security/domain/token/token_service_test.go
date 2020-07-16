package token_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"omics/pkg/security/domain/token"
	"omics/pkg/security/domain/token/mocks"
	cache "omics/pkg/shared/cache/mocks"

	"github.com/golang/mock/gomock"
)

type tokenService struct {
	serv  *token.TokenService
	enc   *mocks.MockTokenEncoder
	cache *cache.MockCache
}

var (
	ErrTest = errors.New("test")
)

func buildTokenService(ctrl *gomock.Controller) *tokenService {
	enc := mocks.NewMockTokenEncoder(ctrl)
	cache := cache.NewMockCache(ctrl)
	serv := token.NewTokenService(cache, enc)

	return &tokenService{
		enc:   enc,
		cache: cache,
		serv:  serv,
	}
}

func dataFixture() token.Data {
	return token.NewData("U002")
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name   string
		mock   func(*tokenService) (context.Context, token.Data)
		rToken token.Token
		rErr   error
	}{{
		"error generating token",
		func(tb *tokenService) (context.Context, token.Data) {
			data := dataFixture()
			tb.enc.EXPECT().
				Encode(gomock.Any()).
				Return(token.Token(""), ErrTest)
			return context.Background(), data
		},
		"",
		ErrTest,
	}, {
		"error saving token",
		func(tb *tokenService) (context.Context, token.Data) {
			data := dataFixture()
			tb.enc.EXPECT().
				Encode(gomock.Any()).
				Return(token.Token("##123##"), nil)
			tb.cache.EXPECT().
				Set(context.Background(), gomock.Any(), data).
				Return(token.Err)
			return context.Background(), data
		},
		"",
		token.Err,
	}, {
		"generate token",
		func(tb *tokenService) (context.Context, token.Data) {
			data := dataFixture()
			tb.enc.EXPECT().
				Encode(gomock.Any()).
				Return(token.Token("##123##"), nil)
			tb.cache.EXPECT().
				Set(context.Background(), gomock.Any(), data).
				Return(nil)
			return context.Background(), data
		},
		"##123##",
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tb := buildTokenService(ctrl)
			aCtx, aData := test.mock(tb)

			rToken, rErr := tb.serv.Create(aCtx, aData)

			if !errors.Is(rErr, test.rErr) {
				t.Errorf("\nExp: %s\nAct: %s", test.rErr, rErr)
			} else if !reflect.DeepEqual(test.rToken, rToken) {
				t.Errorf("\nExp: %s\nAct: %s", test.rToken, rToken)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name  string
		mock  func(*tokenService) (context.Context, token.Token)
		rData token.Data
		rErr  error
	}{{
		"error decoding",
		func(tb *tokenService) (context.Context, token.Token) {
			tb.enc.EXPECT().
				Decode(token.Token("##123##")).
				Return(token.TokenID(""), ErrTest)
			return context.Background(), token.Token("##123##")
		},
		nil,
		ErrTest,
	}, {
		"error getting cache",
		func(tb *tokenService) (context.Context, token.Token) {
			tb.enc.EXPECT().
				Decode(token.Token("##123##")).
				Return(token.TokenID("t123"), nil)
			tb.cache.EXPECT().
				Get(context.Background(), "token:t123").
				Return(nil, ErrTest)
			return context.Background(), token.Token("##123##")
		},
		nil,
		ErrTest,
	}, {
		"validate",
		func(tb *tokenService) (context.Context, token.Token) {
			data := dataFixture()
			tb.enc.EXPECT().
				Decode(token.Token("##123##")).
				Return(token.TokenID("t123"), nil)
			tb.cache.EXPECT().
				Get(context.Background(), "token:t123").
				Return(data, nil)

			return context.Background(), token.Token("##123##")
		},
		dataFixture(),
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tb := buildTokenService(ctrl)
			aCtx, aToken := test.mock(tb)

			rData, rErr := tb.serv.Validate(aCtx, aToken)

			if !errors.Is(rErr, test.rErr) {
				t.Errorf("\nExp: %s\nAct: %s", test.rErr, rErr)
			} else if !reflect.DeepEqual(test.rData, rData) {
				t.Errorf("\nExp: %v\nAct: %v", test.rData, test.rData)
			}
		})
	}
}

func TestInvalidate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		mock func(tb *tokenService) (context.Context, token.Token)
		rErr error
	}{{
		"error decoding",
		func(tb *tokenService) (context.Context, token.Token) {
			tb.enc.EXPECT().
				Decode(token.Token("#123#")).
				Return(token.TokenID(""), ErrTest)
			return context.Background(), token.Token("#123#")
		},
		ErrTest,
	}, {
		"invalidate",
		func(tb *tokenService) (context.Context, token.Token) {
			tb.enc.EXPECT().
				Decode(token.Token("#123#")).
				Return(token.TokenID("123"), nil)
			tb.cache.EXPECT().
				Delete(context.Background(), "token:123").
				Return(nil)
			return context.Background(), token.Token("#123#")
		},
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tb := buildTokenService(ctrl)
			aCtx, aToken := test.mock(tb)

			rErr := tb.serv.Invalidate(aCtx, aToken)

			if !errors.Is(rErr, test.rErr) {
				t.Errorf("\nExp: %s\nAct: %s", test.rErr, rErr)
			}
		})
	}
}
