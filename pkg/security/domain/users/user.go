package users

import (
	"strings"
	"time"

	"omics/pkg/common/models"
)

// User is an aggregate root
type User struct {
	ID        models.ID
	Username  string
	Email     string
	Password  string
	Name      string
	Lastname  string
	Role      Role
	LastLogin time.Time
}

type Role struct {
	Code        string
	Permissions []Permission
}

type Permission struct {
	Permission string
	Module     string
}

func (u *User) IsAdmin() bool {
	return u.Role.Code == "admin"
}

func (u *User) HasPermissions(permission string, module string) bool {
	if u.Role.Code == "admin" {
		return true
	}

	for _, rolePerm := range u.Role.Permissions {
		if rolePerm.Module == module {
			for _, perm := range strings.Split(permission, "") {
				if !strings.Contains(rolePerm.Permission, perm) {
					return false
				}
			}
			return true
		}
	}
	return false
}
