package repository

import (
	"context"

	"todo-api/internal/domain/entity"
)

// TodoRepository defines the interface for todo repository operations
type TodoRepository interface {
	// Create creates a new todo
	Create(ctx context.Context, todo *entity.Todo) error
	
	// GetByID retrieves a todo by its ID
	GetByID(ctx context.Context, id uint) (*entity.Todo, error)
	
	// GetByUserID retrieves todos for a specific user
	GetByUserID(ctx context.Context, userID uint, filter entity.TodoFilter) ([]*entity.Todo, error)
	
	// Update updates a todo
	Update(ctx context.Context, todo *entity.Todo) error
	
	// Delete deletes a todo
	Delete(ctx context.Context, id uint) error
	
	// Count counts todos based on filter
	Count(ctx context.Context, filter entity.TodoFilter) (int64, error)
}
