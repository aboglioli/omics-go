package users

import (
	"context"

	"omics/pkg/security/application/auth"
	"omics/pkg/security/domain/roles"
	"omics/pkg/security/domain/users"
	"omics/pkg/shared/models"
)

type RegisterCommand struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
}

func (cmd *RegisterCommand) Validate() error {
	return nil
}

type UpdateCommand struct {
	Name     string `json:"name"`
	Lastname string `json:"lastnaem"`
}

type ChangePasswordCommand struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type PublicUser struct {
	ID       models.ID `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Lastname string    `json:"lastname"`
}

type UserService interface {
	GetMe(ctx context.Context) (*PublicUser, error)
	GetByID(ctx context.Context, userID models.ID) (*PublicUser, error)

	Register(ctx context.Context, cmd *RegisterCommand) error
	Update(ctx context.Context, userID models.ID, cmd *UpdateCommand) error
	ChangePassword(ctx context.Context, userID string, cmd *ChangePasswordCommand) error
}

type userService struct {
	roleRepo          roles.RoleRepository
	userRepo          users.UserRepository
	passwordHasher    users.PasswordHasher
	authorizationServ auth.AuthorizationService
}

func (s *userService) GetMe(ctx context.Context) (*PublicUser, error) {
	user, err := s.authorizationServ.GetUserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return &PublicUser{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Lastname: user.Lastname,
	}, nil
}

func (s *userService) GetByID(ctx context.Context, userID models.ID) (*users.User, error) {
	user, err := s.authorizationServ.GetUserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if !user.IsAdmin() {
		if !(user.CanRead("users") && user.ID == userID) {
			return nil, ErrUnauthorized
		}
	}

	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrNotFound.Wrap(err)
	}

	return user, nil
}

func (s *userService) Register(ctx context.Context, cmd *RegisterCommand) error {
	if user, err := s.userRepo.FindByUsernameOrEmail(ctx, cmd.Username); user != nil || err == nil {
		if user, err := s.userRepo.FindByUsernameOrEmail(ctx, cmd.Email); user != nil || err == nil {
			return ErrUsers.Code("register").Wrap(err)
		}
	}

	role, err := s.roleRepo.FindByCode(ctx, "user")
	if err != nil {
		return ErrUsers.Code("register").Wrap(err)
	}

	permissions := make([]users.Permission, 0)
	for _, perm := range role.Permissions {
		permissions = append(permissions, users.Permission{
			Permission: perm.Permission,
			Module:     perm.Module.Code,
		})
	}

	user := &users.User{
		ID:       s.userRepo.NextID(),
		Username: cmd.Username,
		Email:    cmd.Email,
		Name:     cmd.Name,
		Lastname: cmd.Lastname,
		Role: users.Role{
			Code:        role.Code,
			Permissions: permissions,
		},
	}

	if err := user.SetPassword(cmd.Password, s.passwordHasher); err != nil {
		return ErrUsers.Code("hash_password").Wrap(err)
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return ErrUsers.Code("register").Wrap(err)
	}

	return nil
}

func (s *userService) Update(ctx context.Context, userID models.ID, cmd *UpdateCommand) error {
	user, err := s.authorizationServ.GetUserFromCtx(ctx)
	if err != nil {
		return err
	}

	if !user.IsAdmin() {
		if !(user.CanUpdate("users") && user.ID == userID) {
			return ErrUnauthorized
		}
	}

	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrNotFound.Wrap(err)
	}

	user.Name = cmd.Name
	user.Lastname = cmd.Lastname

	if err := s.userRepo.Save(ctx, user); err != nil {
		return ErrUsers.Code("update").Wrap(err)
	}

	return nil
}

func (s *userService) ChangePassword(ctx context.Context, userID models.ID, cmd *ChangePasswordCommand) error {
	user, err := s.authorizationServ.GetUserFromCtx(ctx)
	if err != nil {
		return err
	}

	if !user.IsAdmin() {
		if !(user.CanUpdate("users") && user.ID == userID) {
			return ErrUnauthorized
		}
	}

	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrNotFound.Wrap(err)
	}

	if err := user.ChangePassword(cmd.OldPassword, cmd.NewPassword, s.passwordHasher); err != nil {
		return ErrUsers.Code("change_password").AddContext("password", "mismatch").Wrap(err)
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return ErrUsers.Code("change_password").Wrap(err)
	}

	return nil
}
