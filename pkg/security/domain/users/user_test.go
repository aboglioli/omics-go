package users_test

import (
	"strconv"
	"testing"

	"omics/pkg/security/domain/users"
)

func TestPermissions(t *testing.T) {
	tests := []struct {
		aRole           string
		aRolePermission users.Permission
		aPermission     string
		rRes            bool
	}{{
		"user",
		users.Permission{"CR", "module"},
		"CR",
		true,
	}, {
		"user",
		users.Permission{"CR", "module"},
		"C",
		true,
	}, {
		"user",
		users.Permission{"CRUD", "module"},
		"UDC",
		true,
	}, {
		"user",
		users.Permission{"CRD", "module"},
		"U",
		false,
	}, {
		"user",
		users.Permission{"CRD", "module"},
		"CU",
		false,
	}, {
		"user",
		users.Permission{"CD", "module"},
		"DC",
		true,
	}, {
		"user",
		users.Permission{"CRUD", "module"},
		"CRUD",
		true,
	}, {
		"user",
		users.Permission{"CRUD", "module"},
		"CK",
		false,
	}, {
		"admin",
		users.Permission{"R", "module"},
		"CRUD",
		true,
	}}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			user := &users.User{
				Role: users.Role{
					Code:        test.aRole,
					Permissions: []users.Permission{test.aRolePermission},
				},
			}

			rRes := user.IsAdmin() || user.HasPermissions(test.aPermission, "module")
			if rRes != test.rRes {
				t.Errorf("\nExp: %v\nAct: %v", test.rRes, rRes)
			}
		})
	}
}
