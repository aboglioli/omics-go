package token_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	cache "omics/pkg/common/cache/mocks"
	"omics/pkg/security/domain/token"
	"omics/pkg/security/domain/token/mocks"
	"omics/pkg/security/domain/users"

	"github.com/golang/mock/gomock"
)

type testBed struct {
	serv  token.TokenService
	enc   *mocks.MockTokenEncoder
	cache *cache.MockCache
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		mock func(*testBed) (context.Context, *users.User)
		res  token.Token
		err  error
	}{{
		"generate token",
		func(t *testBed) (context.Context, *users.User) {
			user := &users.User{
				ID: "U001",
			}
			t.enc.EXPECT().
				Encode(gomock.Any()).
				Return(token.Token("token:123"), nil)
			t.cache.EXPECT().
				Set(context.Background(), gomock.Any(), user).
				Return(nil)
			return context.Background(), user
		},
		"token:123",
		nil,
	}, {
		"error generating token",
		func(t *testBed) (context.Context, *users.User) {
			user := &users.User{
				ID: "U001",
			}
			t.enc.EXPECT().
				Encode(gomock.Any()).
				Return(token.Token(""), token.ErrToken)
			return context.Background(), user
		},
		"",
		token.ErrToken,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			enc := mocks.NewMockTokenEncoder(ctrl)
			cache := cache.NewMockCache(ctrl)
			serv := token.NewTokenService(cache, enc)

			arg1, arg2 := test.mock(&testBed{
				serv:  serv,
				enc:   enc,
				cache: cache,
			})

			res, err := serv.Create(arg1, arg2)

			if !errors.Is(err, test.err) {
				t.Errorf("\nExp:%v\nAct:%v", test.err, err)
			} else if !reflect.DeepEqual(test.res, res) {
				t.Errorf("\nExp:%v\nAct:%v", test.res, res)
			}
		})
	}

}
