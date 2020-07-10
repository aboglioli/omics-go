package security

import (
	"context"
	"reflect"
	"testing"

	"omics/pkg/models"
)

func buildUserService() *userService {
	// UserRepository
	userRepo := NewInMemUserRepository()
	// TokenService
	enc := FakeTokenEncoder()
	cache := FakeCache()
	tokenServ := &tokenService{
		cache: cache,
		enc:   enc,
	}
	// PasswordHasher
	passwordHasher := FakePasswordHasher()
	return &userService{
		userRepo:       userRepo,
		tokenServ:      tokenServ,
		passwordHasher: passwordHasher,
	}
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
			serv := buildUserService()
			if err := serv.userRepo.Save(context.Background(), test.user); err != nil {
				t.Error(err)
			}
			token, _ := serv.tokenServ.Create(context.Background(), test.user)
			ctx := context.WithValue(context.Background(), "authToken", token.String())

			err := serv.ChangePassword(ctx, test.userID, test.req)

			if test.err != nil {
				if !reflect.DeepEqual(err, test.err) {
					t.Errorf("\nExp:%v\nAct:%v", test.err, err)
				}
			} else {
				savedUser, err := serv.userRepo.FindByID(ctx, test.user.ID)
				if err != nil {
					t.Fatalf("User wasn't saved")
				}

				if savedUser.Password == test.req.NewPassword {
					t.Errorf("Password wasn't hashed")
				}
			}
		})
	}
}
