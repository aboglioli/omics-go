package token_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"omics/pkg/security/domain/token"
	"omics/pkg/security/domain/token/mocks"
	"omics/pkg/security/domain/users"
	cache "omics/pkg/shared/cache/mocks"

	"github.com/golang/mock/gomock"
)

type tokenService struct {
	serv  token.TokenService
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

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name   string
		mock   func(*tokenService) (context.Context, *users.User)
		rToken token.Token
		rErr   error
	}{{
		"error generating token",
		func(tb *tokenService) (context.Context, *users.User) {
			user := &users.User{
				ID: "U001",
			}
			tb.enc.EXPECT().
				Encode(gomock.Any()).
				Return(token.Token(""), ErrTest)
			return context.Background(), user
		},
		"",
		ErrTest,
	}, {
		"error saving token",
		func(tb *tokenService) (context.Context, *users.User) {
			user := &users.User{
				ID: "U001",
			}
			tb.enc.EXPECT().
				Encode(gomock.Any()).
				Return(token.Token("token:123"), nil)
			tb.cache.EXPECT().
				Set(context.Background(), gomock.Any(), user).
				Return(token.ErrToken)
			return context.Background(), user
		},
		"",
		token.ErrToken,
	}, {
		"generate token",
		func(tb *tokenService) (context.Context, *users.User) {
			user := &users.User{
				ID: "U001",
			}
			tb.enc.EXPECT().
				Encode(gomock.Any()).
				Return(token.Token("token:123"), nil)
			tb.cache.EXPECT().
				Set(context.Background(), gomock.Any(), user).
				Return(nil)
			return context.Background(), user
		},
		"token:123",
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tb := buildTokenService(ctrl)
			aCtx, aUser := test.mock(tb)

			rToken, rErr := tb.serv.Create(aCtx, aUser)

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
		rUser *users.User
		rErr  error
	}{{
		"error decoding",
		func(tb *tokenService) (context.Context, token.Token) {
			tb.enc.EXPECT().
				Decode(token.Token("token:123")).
				Return(token.TokenID(""), ErrTest)
			return context.Background(), token.Token("token:123")
		},
		nil,
		ErrTest,
	}, {
		"error getting cache",
		func(tb *tokenService) (context.Context, token.Token) {
			tb.enc.EXPECT().
				Decode(token.Token("token:123")).
				Return(token.TokenID("t123"), nil)
			tb.cache.EXPECT().
				Get(context.Background(), "token:t123").
				Return(nil, ErrTest)
			return context.Background(), token.Token("token:123")
		},
		nil,
		ErrTest,
	}, {
		"validate",
		func(tb *tokenService) (context.Context, token.Token) {
			user := &users.User{
				ID: "#U001",
			}
			tb.enc.EXPECT().
				Decode(token.Token("token:123")).
				Return(token.TokenID("t123"), nil)
			tb.cache.EXPECT().
				Get(context.Background(), "token:t123").
				Return(user, nil)

			return context.Background(), token.Token("token:123")
		},
		&users.User{
			ID: "#U001",
		},
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tb := buildTokenService(ctrl)
			aCtx, aToken := test.mock(tb)

			rUser, rErr := tb.serv.Validate(aCtx, aToken)

			if !errors.Is(rErr, test.rErr) {
				t.Errorf("\nExp: %s\nAct: %s", test.rErr, rErr)
			} else if !reflect.DeepEqual(test.rUser, rUser) {
				t.Errorf("\nExp: %v\nAct: %v", test.rUser, test.rUser)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		mock func(tb *tokenService) (context.Context, token.Token, *users.User)
		rErr error
	}{{
		"error decoding",
		func(tb *tokenService) (context.Context, token.Token, *users.User) {
			tb.enc.EXPECT().
				Decode(token.Token("token:123")).
				Return(token.TokenID(""), ErrTest)
			return context.Background(), token.Token("token:123"), nil
		},
		ErrTest,
	}, {
		"error getting from cache",
		func(tb *tokenService) (context.Context, token.Token, *users.User) {
			tb.enc.EXPECT().
				Decode(token.Token("#123#")).
				Return(token.TokenID("123"), nil)
			tb.cache.EXPECT().
				Get(context.Background(), "token:123").
				Return(nil, ErrTest)
			return context.Background(), token.Token("#123#"), nil
		},
		ErrTest,
	}, {
		"error saving new user",
		func(tb *tokenService) (context.Context, token.Token, *users.User) {
			oldUser := &users.User{
				ID:   "U001",
				Name: "OldName",
			}
			newUser := &users.User{
				ID:   "U001",
				Name: "NewName",
			}
			tb.enc.EXPECT().
				Decode(token.Token("#123#")).
				Return(token.TokenID("123"), nil)
			tb.cache.EXPECT().
				Get(context.Background(), "token:123").
				Return(oldUser, nil)
			tb.cache.EXPECT().
				Set(context.Background(), "token:123", newUser).
				Return(ErrTest)
			return context.Background(), token.Token("#123#"), newUser
		},
		ErrTest,
	}, {
		"update",
		func(tb *tokenService) (context.Context, token.Token, *users.User) {
			oldUser := &users.User{
				ID:   "U001",
				Name: "OldName",
			}
			newUser := &users.User{
				ID:   "U001",
				Name: "NewName",
			}
			tb.enc.EXPECT().
				Decode(token.Token("#123#")).
				Return(token.TokenID("123"), nil)
			tb.cache.EXPECT().
				Get(context.Background(), "token:123").
				Return(oldUser, nil)
			tb.cache.EXPECT().
				Set(context.Background(), "token:123", newUser).
				Return(nil)
			return context.Background(), token.Token("#123#"), newUser
		},
		nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tb := buildTokenService(ctrl)
			aCtx, aToken, aUser := test.mock(tb)

			rErr := tb.serv.Update(aCtx, aToken, aUser)

			if !errors.Is(rErr, test.rErr) {
				t.Errorf("\nExp: %s\nAct: %s", test.rErr, rErr)
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
