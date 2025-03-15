package usecase

import (
	"context"
	"errors"
	"strings"

	"todo-api/internal/domain/entity"
	"todo-api/internal/domain/repository"
	"todo-api/internal/util/jwt"
	"todo-api/internal/util/password"
)

// Errors related to user operations
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameExists    = errors.New("username already exists")
	ErrEmailExists       = errors.New("email already exists")
	ErrInvalidUserData   = errors.New("invalid user data")
	ErrUserCreateFailed  = errors.New("failed to create user")
	ErrUserUpdateFailed  = errors.New("failed to update user")
)

// LoginResponse represents the response of a successful login
type LoginResponse struct {
	Token string
	User  *entity.User
}

// UserUseCase defines the interface for user use cases
type UserUseCase interface {
	Register(ctx context.Context, username, email, password string) (*entity.User, error)
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	UpdateProfile(ctx context.Context, id uint, username, email string) (*entity.User, error)
	UpdatePassword(ctx context.Context, id uint, currentPassword, newPassword string) error
}

// userUseCase implements UserUseCase
type userUseCase struct {
	userRepo repository.UserRepository
	jwt      jwt.JWTService
}

// NewUserUseCase creates a new UserUseCase
func NewUserUseCase(userRepo repository.UserRepository, jwtService jwt.JWTService) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		jwt:      jwtService,
	}
}

// Register registers a new user
func (uc *userUseCase) Register(ctx context.Context, username, email, pwd string) (*entity.User, error) {
	// Validate input
	if username == "" || email == "" || pwd == "" {
		return nil, ErrInvalidUserData
	}

	// Normalize email and username
	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.TrimSpace(username)

	// Check if username exists
	exists, err := uc.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExists
	}

	// Check if email exists
	exists, err = uc.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := password.Hash(pwd)
	if err != nil {
		return nil, err
	}

	// Create user
	user := entity.NewUser(username, email, hashedPassword)
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, ErrUserCreateFailed
	}

	return user, nil
}

// Login logs in a user
func (uc *userUseCase) Login(ctx context.Context, email, pwd string) (*LoginResponse, error) {
	// Normalize email
	email = strings.ToLower(strings.TrimSpace(email))

	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if !password.Verify(pwd, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := uc.jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetUserByID retrieves a user by their ID
func (uc *userUseCase) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateProfile updates a user's profile
func (uc *userUseCase) UpdateProfile(ctx context.Context, id uint, username, email string) (*entity.User, error) {
	// Validate input
	if username == "" || email == "" {
		return nil, ErrInvalidUserData
	}

	// Normalize email and username
	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.TrimSpace(username)

	// Get the user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if username is already taken (if changed)
	if username != user.Username {
		exists, err := uc.userRepo.ExistsByUsername(ctx, username)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrUsernameExists
		}
	}

	// Check if email is already taken (if changed)
	if email != user.Email {
		exists, err := uc.userRepo.ExistsByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEmailExists
		}
	}

	// Update profile
	user.UpdateProfile(username, email)
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, ErrUserUpdateFailed
	}

	return user, nil
}

// UpdatePassword updates a user's password
func (uc *userUseCase) UpdatePassword(ctx context.Context, id uint, currentPassword, newPassword string) error {
	// Validate input
	if currentPassword == "" || newPassword == "" {
		return ErrInvalidUserData
	}

	// Get the user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify current password
	if !password.Verify(currentPassword, user.Password) {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := password.Hash(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.UpdatePassword(hashedPassword)
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return ErrUserUpdateFailed
	}

	return nil
}
