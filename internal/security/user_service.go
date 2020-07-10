package security

import (
	"context"

	"omics/pkg/errors"
	"omics/pkg/models"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
}

func (req *RegisterRequest) Validate() error {
	return nil
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username"`
	Password        string `json:"password"`
}

type LoginResponse struct {
	AuthToken string `json:"auth_token"`
}

type UpdateRequest struct {
	Name     string `json:"name"`
	Lastname string `json:"lastnaem"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UserService interface {
	GetByID(ctx context.Context, userID models.ID) (*User, error)
	GetLoggedIn(ctx context.Context) (*User, error)

	Register(ctx context.Context, req *RegisterRequest) error
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Update(ctx context.Context, userID models.ID, req *UpdateRequest) error
	ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error
	Logout(ctx context.Context)
}

type userService struct {
	roleRepo       RoleRepository
	userRepo       UserRepository
	tokenServ      TokenService
	passwordHasher PasswordHasher
}

func (s *userService) GetByID(ctx context.Context, userID models.ID) (*User, error) {
	token, err := TokenFromContext(ctx)
	if err != nil {
		return nil, errors.ErrTODO
	}

	user, err := s.tokenServ.Validate(ctx, token)
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

func (s *userService) GetLoggedIn(ctx context.Context) (*User, error) {
	token, err := TokenFromContext(ctx)
	if err != nil {
		return nil, errors.ErrTODO
	}

	user, err := s.tokenServ.Validate(ctx, token)
	if err != nil {
		return nil, errors.ErrTODO
	}

	return user, nil
}

func (s *userService) Register(ctx context.Context, req *RegisterRequest) error {
	if user, err := s.userRepo.FindByUsernameOrEmail(ctx, req.Username); user != nil || err == nil {
		if user, err := s.userRepo.FindByUsernameOrEmail(ctx, req.Email); user != nil || err == nil {
			return errors.ErrTODO
		}
	}

	role, err := s.roleRepo.FindByCode(ctx, "user")
	if err != nil {
		return errors.ErrTODO
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Name:     req.Name,
		Lastname: req.Lastname,
		Role:     role,
	}

	hashedPassword, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		return errors.ErrTODO
	}

	user.Password = hashedPassword

	if err := s.userRepo.Save(ctx, user); err != nil {
		return errors.ErrTODO
	}

	return nil
}

func (s *userService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsernameOrEmail(ctx, req.UsernameOrEmail)
	if err != nil {
		return nil, errors.ErrTODO
	}

	if !s.passwordHasher.Compare(user.Password, req.Password) {
		return nil, errors.ErrTODO
	}

	token, err := s.tokenServ.Create(ctx, user)
	if err != nil {
		return nil, errors.ErrTODO
	}

	return &LoginResponse{
		AuthToken: string(token),
	}, nil
}

func (s *userService) Update(ctx context.Context, userID models.ID, req *UpdateRequest) error {
	token, err := TokenFromContext(ctx)
	if err != nil {
		return errors.ErrTODO
	}

	user, err := s.tokenServ.Validate(ctx, token)
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

	user.Name = req.Name
	user.Lastname = req.Lastname

	if err := s.userRepo.Save(ctx, user); err != nil {
		return errors.ErrTODO
	}

	return nil
}

func (s *userService) ChangePassword(ctx context.Context, userID models.ID, req *ChangePasswordRequest) error {
	token, err := TokenFromContext(ctx)
	if err != nil {
		return errors.ErrTODO
	}

	user, err := s.tokenServ.Validate(ctx, token)
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

	if !s.passwordHasher.Compare(user.Password, req.OldPassword) {
		return errors.ErrTODO
	}

	hashedPassword, err := s.passwordHasher.Hash(req.NewPassword)
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
	token, err := TokenFromContext(ctx)
	if err != nil {
		return errors.ErrTODO
	}

	if err := s.tokenServ.Invalidate(ctx, token); err != nil {
		return errors.ErrTODO
	}

	return nil
}
