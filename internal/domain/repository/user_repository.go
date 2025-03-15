package repository

import (
	"context"

	"todo-api/internal/domain/entity"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error
	
	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	
	// GetByUsername retrieves a user by their username
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	
	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	
	// Update updates a user's information
	Update(ctx context.Context, user *entity.User) error
	
	// Delete deletes a user
	Delete(ctx context.Context, id uint) error
	
	// ExistsByUsername checks if a user with the given username exists
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	
	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
