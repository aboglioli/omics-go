package users

import (
	"strconv"
	"testing"
)

func TestPermissions(t *testing.T) {
	tests := []struct {
		role           string
		rolePermission Permission
		permission     string
		res            bool
	}{{
		"user",
		Permission{"CR", "module"},
		"CR",
		true,
	}, {
		"user",
		Permission{"CR", "module"},
		"C",
		true,
	}, {
		"user",
		Permission{"CRUD", "module"},
		"UDC",
		true,
	}, {
		"user",
		Permission{"CRD", "module"},
		"U",
		false,
	}, {
		"user",
		Permission{"CRD", "module"},
		"CU",
		false,
	}, {
		"user",
		Permission{"CD", "module"},
		"DC",
		true,
	}, {
		"user",
		Permission{"CRUD", "module"},
		"CRUD",
		true,
	}, {
		"admin",
		Permission{"R", "module"},
		"CRUD",
		true,
	}}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			user := &User{
				Role: Role{
					Code:        test.role,
					Permissions: []Permission{test.rolePermission},
				},
			}

			res := user.HasPermissions(test.permission, "module")
			if res != test.res {
				t.Errorf("\nExp:%v\nAct:%v", test.res, res)
			}
		})
	}
}
