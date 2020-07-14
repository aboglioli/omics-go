package users

import (
	"strings"
	"time"

	"omics/pkg/shared/models"
)

// User is an aggregate root
type User struct {
	id                 models.ID
	username           Username
	email              Email
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

	if err := u.SetUsername(username); err != nil {
		errs = errs.Merge(err)
	}

	if err := u.SetEmail(email); err != nil {
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

func (u *User) Username() Username {
	return u.username
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) Name() Fullname {
	return u.name
}

func (u *User) SetUsername(username string) error {
	un, err := NewUsername(username)
	if err != nil {
		return err
	}

	u.username = un
	return nil
}

func (u *User) SetEmail(email string) error {
	e, err := NewEmail(email)
	if err != nil {
		return err
	}

	u.email = e
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

func (u *User) WasAuthenticated() {
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

func (u *User) ComparePassword(plainPassword string, hasher PasswordHasher) bool {
	return hasher.Compare(u.password, plainPassword)
}

func (u *User) ChangePassword(oldPassword, newPassword string, hasher PasswordHasher, validator PasswordValidator) error {
	if u.password != "" {
		if !u.ComparePassword(oldPassword, hasher) {
			return ErrUsers.Code("password_mismatch")
		}
	}

	if err := validator.Validate(newPassword); err != nil {
		return ErrValidation.AddContext("password", "weak").Wrap(err)
	}

	hashedPassword, err := hasher.Hash(newPassword)
	if err != nil {
		return ErrUsers.Code("hash_password").Wrap(err)
	}

	u.password = hashedPassword

	return nil
}
