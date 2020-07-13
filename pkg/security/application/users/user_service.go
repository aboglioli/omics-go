package users

import (
	"context"

	"omics/pkg/common/models"
	"omics/pkg/security/domain/roles"
	"omics/pkg/security/domain/token"
	"omics/pkg/security/domain/users"
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

type LoginCommand struct {
	UsernameOrEmail string `json:"username"`
	Password        string `json:"password"`
}

type LoginResponse struct {
	AuthToken string `json:"auth_token"`
}

type UpdateCommand struct {
	Name     string `json:"name"`
	Lastname string `json:"lastnaem"`
}

type ChangePasswordCommand struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UserService interface {
	GetByID(ctx context.Context, userID models.ID) (*users.User, error)
	GetLoggedIn(ctx context.Context) (*users.User, error)

	Register(ctx context.Context, cmd *RegisterCommand) error
	Login(ctx context.Context, cmd *LoginCommand) (*LoginResponse, error)
	Update(ctx context.Context, userID models.ID, cmd *UpdateCommand) error
	ChangePassword(ctx context.Context, userID string, cmd *ChangePasswordCommand) error
	Logout(ctx context.Context)
}

type userService struct {
	roleRepo       roles.RoleRepository
	userRepo       users.UserRepository
	tokenServ      token.TokenService
	passwordHasher users.PasswordHasher
}

func (s *userService) GetByID(ctx context.Context, userID models.ID) (*users.User, error) {
	t, err := token.TokenFromContext(ctx)
	if err != nil {
		return nil, ErrNotFound.Wrap(err)
	}

	user, err := s.tokenServ.Validate(ctx, t)
	if err != nil {
		return nil, ErrNotFound.Wrap(err)
	}

	if !user.HasPermissions("R", "users") {
		return nil, ErrNotFound.Wrap(err)
	}

	if !user.IsAdmin() && user.ID != userID {
		return nil, ErrNotFound.Wrap(err)
	}

	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrNotFound.Wrap(err)
	}

	return user, nil
}

func (s *userService) GetLoggedIn(ctx context.Context) (*users.User, error) {
	t, err := token.TokenFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized.Wrap(err)
	}

	user, err := s.tokenServ.Validate(ctx, t)
	if err != nil {
		return nil, ErrUnauthorized.Wrap(err)
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

func (s *userService) Login(ctx context.Context, cmd *LoginCommand) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsernameOrEmail(ctx, cmd.UsernameOrEmail)
	if err != nil {
		return nil, ErrUsers.Code("login").Wrap(err)
	}

	if !user.ComparePassword(cmd.Password, s.passwordHasher) {
		return nil, ErrUsers.Code("login").AddContext("password", "mismatch")
	}

	t, err := s.tokenServ.Create(ctx, user)
	if err != nil {
		return nil, ErrUsers.Code("login").Wrap(err)
	}

	return &LoginResponse{
		AuthToken: string(t),
	}, nil
}

func (s *userService) Update(ctx context.Context, userID models.ID, cmd *UpdateCommand) error {
	user, err := s.GetLoggedIn(ctx)
	if err != nil {
		return ErrUsers.Code("update").Wrap(err)
	}

	if !user.HasPermissions("U", "users") {
		return ErrUnauthorized
	}

	if !user.IsAdmin() && user.ID != userID {
		return ErrUnauthorized
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
	user, err := s.GetLoggedIn(ctx)
	if err != nil {
		return ErrUsers.Code("change_password").Wrap(err)
	}

	if !user.HasPermissions("U", "users") {
		return ErrUnauthorized
	}

	if !user.IsAdmin() && user.ID != userID {
		return ErrUnauthorized
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

func (s *userService) Logout(ctx context.Context) error {
	t, err := token.TokenFromContext(ctx)
	if err != nil {
		return ErrUnauthorized.Wrap(err)
	}

	if err := s.tokenServ.Invalidate(ctx, t); err != nil {
		return ErrUnauthorized.Wrap(err)
	}

	return nil
}
