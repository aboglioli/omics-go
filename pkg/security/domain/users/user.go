package users

import (
	"strings"
	"time"

	"omics/pkg/shared/models"
)

// User is an aggregate root
type User struct {
	ID        models.ID
	Username  string
	Email     string
	password  string
	Name      string
	Lastname  string
	Role      Role
	LastLogin time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Role struct {
	Code        string
	Permissions []Permission
}

type Permission struct {
	Permission string
	Module     string
}

func (u *User) HasRole(role string) bool {
	return u.Role.Code == role
}

func (u *User) HasPermissions(permission string, module string) bool {
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

func (u *User) IsAdmin() bool {
	return u.HasRole("admin")
}

func (u *User) CanCreate(module string) bool {
	return u.HasPermissions("C", module)
}

func (u *User) CanRead(module string) bool {
	return u.HasPermissions("R", module)
}

func (u *User) CanUpdate(module string) bool {
	return u.HasPermissions("U", module)
}

func (u *User) CanDelete(module string) bool {
	return u.HasPermissions("D", module)
}

func (u *User) ComparePassword(plainPassword string, hasher PasswordHasher) bool {
	return hasher.Compare(u.password, plainPassword)
}

func (u *User) SetPassword(plainPassword string, hasher PasswordHasher) error {
	hashedPassword, err := hasher.Hash(plainPassword)
	if err != nil {
		return ErrUsers.Code("hash_password").Wrap(err)
	}

	u.password = hashedPassword

	return nil
}

func (u *User) ChangePassword(oldPassword, newPassword string, hasher PasswordHasher) error {
	if !u.ComparePassword(oldPassword, hasher) {
		return ErrUsers.Code("password_mismatch")
	}

	u.SetPassword(newPassword, hasher)

	return nil
}
