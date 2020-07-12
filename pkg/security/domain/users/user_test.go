package users_test

import (
	"strconv"
	"testing"

	"omics/pkg/security/domain/users"
)

func TestPermissions(t *testing.T) {
	tests := []struct {
		role           string
		rolePermission users.Permission
		permission     string
		res            bool
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
		"admin",
		users.Permission{"R", "module"},
		"CRUD",
		true,
	}}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			user := &users.User{
				Role: users.Role{
					Code:        test.role,
					Permissions: []users.Permission{test.rolePermission},
				},
			}

			res := user.HasPermissions(test.permission, "module")
			if res != test.res {
				t.Errorf("\nExp:%v\nAct:%v", test.res, res)
			}
		})
	}
}
