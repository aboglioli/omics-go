package security

import (
	"context"
	"fmt"
	"omics/pkg/models"
	"reflect"
	"testing"
)

type fakePasswordHasher struct{}

func (ph *fakePasswordHasher) Hash(plainPassword string) (string, error) {
	return fmt.Sprintf("#%s#", plainPassword), nil
}

func (ph *fakePasswordHasher) Compare(hashedPassword, plainPassword string) bool {
	return hashedPassword == fmt.Sprintf("#%s#", plainPassword)
}

type fakeTokenService struct {
	user *User
}

func (s *fakeTokenService) Create(ctx context.Context, user *User) (Token, error) {
	return "token", nil
}

func (s *fakeTokenService) Validate(ctx context.Context, token Token) (*User, error) {
	if token == "token" {
		return s.user, nil
	}
	return nil, ErrNull
}

func (s *fakeTokenService) ValidateFromContext(ctx context.Context) (*User, error) {
	token, ok := ctx.Value("authToken").(string)
	if !ok {
		return nil, ErrNull
	}
	return s.Validate(ctx, Token(token))
}

func (s *fakeTokenService) Invalidate(ctx context.Context, token Token) error {
	return nil
}

func (s *fakeTokenService) InvalidateFromContext(ctx context.Context) error {
	return nil
}

func TestChangePassword(t *testing.T) {
	tests := []struct {
		name   string
		user   *User
		userID models.ID
		req    *ChangePasswordRequest
		err    error
	}{
		{
			"valid user",
			&User{
				ID:       1,
				Password: "#123#",
				Role: Role{
					Code: "user",
					Permissions: []Permission{
						Permission{"CRUD", "users"},
					},
				},
			},
			1,
			&ChangePasswordRequest{
				OldPassword: "123",
				NewPassword: "456",
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userRepo := &inmemUserRepository{
				users: []*User{test.user},
			}
			passwordHasher := &fakePasswordHasher{}
			serv := userService{
				userRepo:       userRepo,
				passwordHasher: passwordHasher,
				tokenServ:      &fakeTokenService{test.user},
			}
			ctx := context.WithValue(context.Background(), "authToken", "token")

			err := serv.ChangePassword(ctx, test.userID, test.req)

			if test.err != nil {
				if !reflect.DeepEqual(err, test.err) {
					t.Errorf("\nExp:%v\nAct:%v", test.err, err)
				}
			} else {
				changedPassword := userRepo.users[0].Password
				shouldBe, _ := passwordHasher.Hash(test.req.NewPassword)
				if changedPassword == test.req.NewPassword || changedPassword != shouldBe {
					t.Errorf("Password wasn't hashed")
				}
			}
		})
	}
}
