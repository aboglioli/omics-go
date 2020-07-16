package users

import (
	"context"

	"omics/pkg/security/domain/roles"
	"omics/pkg/security/domain/token"
	"omics/pkg/security/domain/users"
	"omics/pkg/shared/models"
)

type UserService struct {
	authenticationServ *users.AuthenticationService
	authorizationServ  *users.AuthorizationService
	roleRepo           roles.RoleRepository
	tokenServ          *token.TokenService
	userRepo           users.UserRepository
	userServ           *users.UserService
	validationRepo     users.ValidationRepository
}

func NewUserService(
	authenticationServ *users.AuthenticationService,
	authorizationServ *users.AuthorizationService,
	roleRepo roles.RoleRepository,
	tokenServ *token.TokenService,
	userRepo users.UserRepository,
	userServ *users.UserService,
	validationRepo users.ValidationRepository,
) *UserService {
	return &UserService{
		authenticationServ: authenticationServ,
		authorizationServ:  authorizationServ,
		roleRepo:           roleRepo,
		tokenServ:          tokenServ,
		userRepo:           userRepo,
		userServ:           userServ,
		validationRepo:     validationRepo,
	}
}

func (s *UserService) Me(ctx context.Context) (*users.User, error) {
	t, err := token.FromContext(ctx)
	if err != nil {
		return nil, users.ErrUnauthorized.Wrap(err)
	}

	data, err := s.tokenServ.Validate(ctx, t)
	if err != nil {
		return nil, users.ErrUnauthorized.Wrap(err)
	}

	userID, err := data.UserID()
	if err != nil {
		return nil, users.ErrUnauthorized.Wrap(err)
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, users.ErrNotFound.Wrap(err)
	}

	if !user.IsActive() {
		return nil, users.ErrUnauthorized
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, userID models.ID) (*users.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, users.ErrNotFound.Wrap(err)
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, cmd *LoginCommand) (token.Token, error) {
	user, err := s.authenticationServ.Authenticate(ctx, cmd.UsernameOrEmail, cmd.Password)
	if err != nil {
		return "", err
	}

	data := token.NewData(user.ID())
	tok, err := s.tokenServ.Create(ctx, data)
	if err != nil {
		return "", users.ErrUnauthorized.Wrap(err)
	}

	return tok, nil
}

func (s *UserService) Register(ctx context.Context, cmd *RegisterCommand) error {
	if err := cmd.Validate(); err != nil {
		return users.ErrValidation.Merge(err)
	}

	if err := s.userServ.Available(ctx, cmd.Username, cmd.Email); err != nil {
		return err
	}

	role, err := s.roleRepo.FindByCode(ctx, roles.USER)
	if err != nil {
		return roles.ErrNotFound.Wrap(err)
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
	if err := s.userServ.ChangePassword(user, "", cmd.Password); err != nil {
		return err
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return users.Err.Wrap(err)
	}

	v := users.NewValidation(user.ID())
	if err := s.validationRepo.Save(ctx, v); err != nil {
		return users.Err.Wrap(err)
	}

	return nil
}

func (s *UserService) Update(ctx context.Context, userID models.ID, cmd *UpdateCommand) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return users.ErrNotFound.Wrap(err)
	}

	user.SetName(cmd.Name, cmd.Lastname)

	if err := s.userRepo.Save(ctx, user); err != nil {
		return users.Err.Wrap(err)
	}

	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID models.ID, cmd *ChangePasswordCommand) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return users.ErrNotFound.Wrap(err)
	}

	if !user.IsActive() {
		return users.ErrUnauthorized.AddContext("active", "false")
	}

	if err := s.userServ.ChangePassword(user, cmd.OldPassword, cmd.NewPassword); err != nil {
		return err
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return users.Err.Wrap(err)
	}

	return nil
}

func (s *UserService) Validate(ctx context.Context, userID models.ID, code string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return users.ErrNotFound.Wrap(err)
	}

	v, err := s.validationRepo.FindByUserID(ctx, userID)
	if err != nil {
		return users.ErrNotFound.AddContext("validation", "not_found").Wrap(err)
	}

	if err := v.Validate(user, code); err != nil {
		return err
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return users.Err.Wrap(err)
	}

	if err := s.validationRepo.Delete(ctx, userID); err != nil {
		return users.Err.Wrap(err)
	}

	return nil
}
