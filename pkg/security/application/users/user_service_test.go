package users

import (
	"context"
	"reflect"
	"testing"

	"omics/pkg/common/models"
	"omics/pkg/security/domain/token"
	"omics/pkg/security/domain/users"
	"omics/pkg/security/infrastructure/persistence"
	"omics/pkg/security/mocks"
)

func buildUserService() *userService {
	// UserRepository
	userRepo := persistence.NewInMemUserRepository()
	// TokenService
	enc := mocks.FakeTokenEncoder()
	cache := mocks.FakeCache()
	tokenServ := token.NewTokenService(cache, enc)
	// PasswordHasher
	passwordHasher := mocks.FakePasswordHasher()
	return &userService{
		userRepo:       userRepo,
		tokenServ:      tokenServ,
		passwordHasher: passwordHasher,
	}
}

func TestChangePassword(t *testing.T) {
	tests := []struct {
		name   string
		user   *users.User
		userID models.ID
		req    *ChangePasswordCommand
		err    error
	}{
		{
			"valid user",
			&users.User{
				ID:       "U01",
				Password: "#123#",
				Role: users.Role{
					Code: "user",
					Permissions: []users.Permission{
						users.Permission{"CRUD", "users"},
					},
				},
			},
			"U01",
			&ChangePasswordCommand{
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
