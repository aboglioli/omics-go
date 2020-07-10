package users

import (
	"context"

	"omics/pkg/common/errors"
	"omics/pkg/common/models"
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
	roleRepo       users.RoleRepository
	userRepo       users.UserRepository
	tokenServ      token.TokenService
	passwordHasher users.PasswordHasher
}

func (s *userService) GetByID(ctx context.Context, userID models.ID) (*users.User, error) {
	t, err := token.TokenFromContext(ctx)
	if err != nil {
		return nil, errors.ErrTODO
	}

	user, err := s.tokenServ.Validate(ctx, t)
	if err != nil {
		return nil, errors.ErrTODO
	}

	if !user.HasPermissions("R", "users") {
		return nil, errors.ErrTODO
	}

	if !user.IsAdmin() && user.ID != userID {
		return nil, errors.ErrTODO
	}

	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrTODO
	}

	return user, nil
}

func (s *userService) GetLoggedIn(ctx context.Context) (*users.User, error) {
	t, err := token.TokenFromContext(ctx)
	if err != nil {
		return nil, errors.ErrTODO
	}

	user, err := s.tokenServ.Validate(ctx, t)
	if err != nil {
		return nil, errors.ErrTODO
	}

	return user, nil
}

func (s *userService) Register(ctx context.Context, cmd *RegisterCommand) error {
	if user, err := s.userRepo.FindByUsernameOrEmail(ctx, cmd.Username); user != nil || err == nil {
		if user, err := s.userRepo.FindByUsernameOrEmail(ctx, cmd.Email); user != nil || err == nil {
			return errors.ErrTODO
		}
	}

	role, err := s.roleRepo.FindByCode(ctx, "user")
	if err != nil {
		return errors.ErrTODO
	}

	user := &users.User{
		Username: cmd.Username,
		Email:    cmd.Email,
		Name:     cmd.Name,
		Lastname: cmd.Lastname,
		Role:     role,
	}

	hashedPassword, err := s.passwordHasher.Hash(cmd.Password)
	if err != nil {
		return errors.ErrTODO
	}

	user.Password = hashedPassword

	if err := s.userRepo.Save(ctx, user); err != nil {
		return errors.ErrTODO
	}

	return nil
}

func (s *userService) Login(ctx context.Context, cmd *LoginCommand) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsernameOrEmail(ctx, cmd.UsernameOrEmail)
	if err != nil {
		return nil, errors.ErrTODO
	}

	if !s.passwordHasher.Compare(user.Password, cmd.Password) {
		return nil, errors.ErrTODO
	}

	t, err := s.tokenServ.Create(ctx, user)
	if err != nil {
		return nil, errors.ErrTODO
	}

	return &LoginResponse{
		AuthToken: string(t),
	}, nil
}

func (s *userService) Update(ctx context.Context, userID models.ID, cmd *UpdateCommand) error {
	t, err := token.TokenFromContext(ctx)
	if err != nil {
		return errors.ErrTODO
	}

	user, err := s.tokenServ.Validate(ctx, t)
	if err != nil {
		return errors.ErrTODO
	}

	if !user.HasPermissions("U", "users") {
		return errors.ErrTODO
	}

	if !user.IsAdmin() && user.ID != userID {
		return errors.ErrTODO
	}

	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.ErrTODO
	}

	user.Name = cmd.Name
	user.Lastname = cmd.Lastname

	if err := s.userRepo.Save(ctx, user); err != nil {
		return errors.ErrTODO
	}

	return nil
}

func (s *userService) ChangePassword(ctx context.Context, userID models.ID, cmd *ChangePasswordCommand) error {
	t, err := token.TokenFromContext(ctx)
	if err != nil {
		return errors.ErrTODO
	}

	user, err := s.tokenServ.Validate(ctx, t)
	if err != nil {
		return errors.ErrTODO
	}

	if !user.HasPermissions("U", "users") {
		return errors.ErrTODO
	}

	if !user.IsAdmin() && user.ID != userID {
		return errors.ErrTODO
	}

	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.ErrTODO
	}

	if !s.passwordHasher.Compare(user.Password, cmd.OldPassword) {
		return errors.ErrTODO
	}

	hashedPassword, err := s.passwordHasher.Hash(cmd.NewPassword)
	if err != nil {
		return errors.ErrTODO
	}

	user.Password = hashedPassword

	if err := s.userRepo.Save(ctx, user); err != nil {
		return errors.ErrTODO
	}

	return nil
}

func (s *userService) Logout(ctx context.Context) error {
	t, err := token.TokenFromContext(ctx)
	if err != nil {
		return errors.ErrTODO
	}

	if err := s.tokenServ.Invalidate(ctx, t); err != nil {
		return errors.ErrTODO
	}

	return nil
}
