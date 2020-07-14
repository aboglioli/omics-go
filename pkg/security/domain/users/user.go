package users

import (
	"strings"
	"time"

	"omics/pkg/shared/models"
)

// User is an aggregate root
type User struct {
	id                 models.ID
	identity           Identity
	password           string
	name               Fullname
	role               Role
	lastAuthentication time.Time
}

func NewUser(
	id models.ID,
	username string,
	email string,
	name string,
	lastname string,
) (*User, error) {
	errs := ErrValidation

	u := &User{id: id}

	if err := u.setIdentity(username, email); err != nil {
		errs = errs.Merge(err)
	}

	if err := u.SetName(name, lastname); err != nil {
		errs = errs.Merge(err)
	}

	if errs.ContextLen() > 0 {
		return nil, errs
	}

	return u, nil
}

func (u *User) ID() models.ID {
	return u.id
}

func (u *User) Identity() Identity {
	return u.identity
}

func (u *User) Name() Fullname {
	return u.name
}

func (u *User) setIdentity(username, email string) error {
	id, err := NewIdentity(username, email)
	if err != nil {
		return err
	}

	u.identity = id
	return nil
}

func (u *User) SetName(name, lastname string) error {
	n, err := NewFullname(name, lastname)
	if err != nil {
		return err
	}

	u.name = n
	return nil
}

func (u *User) AssignRole(role Role) error {
	u.role = role

	return nil
}

func (u *User) wasAuthenticated() {
	u.lastAuthentication = time.Now()
}

func (u *User) Role() Role {
	return u.role
}

func (u *User) HasRole(role string) bool {
	return u.role.Code == role
}

func (u *User) HasPermissions(permissions string, module string) bool {
	for _, rolePerm := range u.role.Permissions {
		if rolePerm.Module == module {
			for _, perm := range strings.Split(permissions, "") {
				if !strings.Contains(rolePerm.Permission, perm) {
					return false
				}
			}
			return true
		}
	}
	return false
}
