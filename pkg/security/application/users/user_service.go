package users

import (
	"context"

	"omics/pkg/security/application/auth"
	"omics/pkg/security/domain/token"
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

type UserService interface {
	Me(ctx context.Context) (*users.User, error)
	GetByID(ctx context.Context, userID models.ID) (*users.User, error)

	Register(ctx context.Context, cmd *RegisterCommand) error
	Update(ctx context.Context, userID models.ID, cmd *UpdateCommand) error
	ChangePassword(ctx context.Context, userID string, cmd *ChangePasswordCommand) error
}

type userService struct {
	roleRepo          users.RoleRepository
	userRepo          users.UserRepository
	userServ          users.UserService
	passwordHasher    users.PasswordHasher
	passwordValidator users.PasswordValidator
	authorizationServ auth.AuthorizationService
}

func (s *userService) Me(ctx context.Context) (*users.User, error) {
	t, err := token.FromContext(ctx)
	if err != nil {
		return nil, ErrUsers.Wrap(err)
	}

	user, err := s.authorizationServ.GetUserByToken(ctx, t)
	if err != nil {
		return nil, ErrUsers.Wrap(err)
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, userID models.ID) (*users.User, error) {
	user, err := s.Me(ctx)
	if err != nil {
		return nil, ErrUsers.Wrap(err)
	}

	if !user.HasRole(users.ADMIN) {
		if !(user.HasPermissions(users.READ, "users") && user.ID().Equals(userID)) {
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
	if err := cmd.Validate(); err != nil {
		return ErrUsers.Code("register").Wrap(err)
	}

	if err := s.userServ.Available(ctx, cmd.Username, cmd.Email); err != nil {
		return ErrUsers.Code("register").Wrap(err)
	}

	role, err := s.roleRepo.FindByCode(ctx, users.USER)
	if err != nil {
		return ErrUsers.Code("register").Wrap(err)
	}

	user, err := users.NewUser(
		s.userRepo.NextID(),
		cmd.Username,
		cmd.Email,
		cmd.Name,
		cmd.Lastname,
	)

	if err != nil {
		return err
	}

	user.AssignRole(role)

	if err := user.ChangePassword("", cmd.Password, s.passwordHasher, s.passwordValidator); err != nil {
		return ErrUsers.Code("hash_password").Wrap(err)
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return ErrUsers.Code("register").Wrap(err)
	}

	return nil
}

func (s *userService) Update(ctx context.Context, userID models.ID, cmd *UpdateCommand) error {
	user, err := s.Me(ctx)
	if err != nil {
		return ErrUsers.Wrap(err)
	}

	if !user.HasRole(users.ADMIN) {
		if !(user.HasPermissions(users.UPDATE, "users") && user.ID().Equals(userID)) {
			return ErrUnauthorized
		}
	}

	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrNotFound.Wrap(err)
	}

	user.SetName(cmd.Name, cmd.Lastname)

	if err := s.userRepo.Save(ctx, user); err != nil {
		return ErrUsers.Code("update").Wrap(err)
	}

	return nil
}

func (s *userService) ChangePassword(ctx context.Context, userID models.ID, cmd *ChangePasswordCommand) error {
	user, err := s.Me(ctx)
	if err != nil {
		return ErrUsers.Wrap(err)
	}

	if !user.HasRole(users.ADMIN) {
		if !(user.HasPermissions(users.UPDATE, "users") && user.ID().Equals(userID)) {
			return ErrUnauthorized
		}
	}

	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrNotFound.Wrap(err)
	}

	if err := user.ChangePassword(cmd.OldPassword, cmd.NewPassword, s.passwordHasher, s.passwordValidator); err != nil {
		return ErrUsers.Code("change_password").AddContext("password", "mismatch").Wrap(err)
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return ErrUsers.Code("change_password").Wrap(err)
	}

	return nil
}
