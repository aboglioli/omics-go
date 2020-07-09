package security

import (
	"context"
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
	userRepo       UserRepository
	roleRepo       RoleRepository
	tokenServ      TokenService
	passwordHasher PasswordHasher
}

func (s *userService) GetByID(ctx context.Context, userID models.ID) (*User, error) {
	user, err := s.tokenServ.ValidateFromContext(ctx)
	if err != nil {
		return nil, ErrNull
	}

	if !user.HasPermissions("R", "users") {
		return nil, ErrNull
	}

	if !user.IsAdmin() && user.ID != userID {
		return nil, ErrNull
	}

	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrNull
	}

	return user, nil
}

func (s *userService) GetLoggedIn(ctx context.Context) (*User, error) {
	user, err := s.tokenServ.ValidateFromContext(ctx)
	if err != nil {
		return nil, ErrNull
	}

	return user, nil
}

func (s *userService) Register(ctx context.Context, req *RegisterRequest) error {
	if user, err := s.userRepo.FindByUsernameOrEmail(ctx, req.Username); user != nil || err == nil {
		if user, err := s.userRepo.FindByUsernameOrEmail(ctx, req.Email); user != nil || err == nil {
			return ErrNull
		}
	}

	role, err := s.roleRepo.FindByCode(ctx, "user")
	if err != nil {
		return ErrNull
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
		return ErrNull
	}

	user.Password = hashedPassword

	if err := s.userRepo.Save(ctx, user); err != nil {
		return ErrNull
	}

	return nil
}

func (s *userService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsernameOrEmail(ctx, req.UsernameOrEmail)
	if err != nil {
		return nil, ErrNull
	}

	if !s.passwordHasher.Compare(user.Password, req.Password) {
		return nil, ErrNull
	}

	token, err := s.tokenServ.Create(ctx, user)
	if err != nil {
		return nil, ErrNull
	}

	return &LoginResponse{
		AuthToken: string(token),
	}, nil
}

func (s *userService) Update(ctx context.Context, userID models.ID, req *UpdateRequest) error {
	user, err := s.tokenServ.ValidateFromContext(ctx)
	if err != nil {
		return ErrNull
	}

	if !user.HasPermissions("U", "users") {
		return ErrNull
	}

	if !user.IsAdmin() && user.ID != userID {
		return ErrNull
	}

	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrNull
	}

	user.Name = req.Name
	user.Lastname = req.Lastname

	if err := s.userRepo.Save(ctx, user); err != nil {
		return ErrNull
	}

	return nil
}

func (s *userService) ChangePassword(ctx context.Context, userID models.ID, req *ChangePasswordRequest) error {
	user, err := s.tokenServ.ValidateFromContext(ctx)
	if err != nil {
		return ErrNull
	}

	if !user.HasPermissions("U", "users") {
		return ErrNull
	}

	if !user.IsAdmin() && user.ID != userID {
		return ErrNull
	}

	user, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrNull
	}

	if !s.passwordHasher.Compare(user.Password, req.OldPassword) {
		return ErrNull
	}

	hashedPassword, err := s.passwordHasher.Hash(req.NewPassword)
	if err != nil {
		return ErrNull
	}

	user.Password = hashedPassword

	if err := s.userRepo.Save(ctx, user); err != nil {
		return ErrNull
	}

	return nil
}

func (s *userService) Logout(ctx context.Context) error {
	if err := s.tokenServ.InvalidateFromContext(ctx); err != nil {
		return ErrNull
	}
	return nil
}
