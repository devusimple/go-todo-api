package postgres

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"todo-api/internal/domain/entity"
	"todo-api/internal/domain/repository"
)

// todoRepository implements repository.TodoRepository
type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates a new TodoRepository
func NewTodoRepository(db *gorm.DB) repository.TodoRepository {
	return &todoRepository{
		db: db,
	}
}

// Create creates a new todo
func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

// GetByID retrieves a todo by its ID
func (r *todoRepository) GetByID(ctx context.Context, id uint) (*entity.Todo, error) {
	var todo entity.Todo
	err := r.db.WithContext(ctx).First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// GetByUserID retrieves todos for a specific user
func (r *todoRepository) GetByUserID(ctx context.Context, userID uint, filter entity.TodoFilter) ([]*entity.Todo, error) {
	var todos []*entity.Todo
	
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	
	// Apply filters
	if filter.Completed != nil {
		query = query.Where("completed = ?", *filter.Completed)
	}
	
	// Apply search if provided
	if filter.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(filter.Search))
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", searchQuery, searchQuery)
	}
	
	// Apply pagination
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)
	
	// Apply ordering
	query = query.Order("created_at DESC")
	
	err := query.Find(&todos).Error
	if err != nil {
		return nil, err
	}
	
	return todos, nil
}

// Update updates a todo
func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	return r.db.WithContext(ctx).Save(todo).Error
}

// Delete deletes a todo
func (r *todoRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Todo{}, id).Error
}

// Count counts todos based on filter
func (r *todoRepository) Count(ctx context.Context, filter entity.TodoFilter) (int64, error) {
	var count int64
	
	query := r.db.WithContext(ctx).Model(&entity.Todo{}).Where("user_id = ?", filter.UserID)
	
	// Apply filters
	if filter.Completed != nil {
		query = query.Where("completed = ?", *filter.Completed)
	}
	
	// Apply search if provided
	if filter.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(filter.Search))
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", searchQuery, searchQuery)
	}
	
	err := query.Count(&count).Error
	return count, err
}
