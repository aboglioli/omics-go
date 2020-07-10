package security

import (
	"context"
	"reflect"
	"testing"

	"omics/pkg/models"
)

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
			passwordHasher := FakePasswordHasher()
			tokenServ := FakeTokenService()
			token, _ := tokenServ.Create(context.Background(), test.user)
			serv := userService{
				userRepo:       userRepo,
				passwordHasher: passwordHasher,
				tokenServ:      tokenServ,
			}
			ctx := context.WithValue(context.Background(), "authToken", token.String())

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
