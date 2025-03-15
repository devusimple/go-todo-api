package usecase

import (
	"context"
	"errors"
	"time"

	"todo-api/internal/domain/entity"
	"todo-api/internal/domain/repository"
)

// Errors related to todo operations
var (
	ErrTodoNotFound     = errors.New("todo not found")
	ErrNotAuthorized    = errors.New("not authorized to access this todo")
	ErrInvalidTodoData  = errors.New("invalid todo data")
	ErrTodoCreateFailed = errors.New("failed to create todo")
	ErrTodoUpdateFailed = errors.New("failed to update todo")
	ErrTodoDeleteFailed = errors.New("failed to delete todo")
)

// TodoUseCase defines the interface for todo use cases
type TodoUseCase interface {
	CreateTodo(ctx context.Context, title, description string, userID uint) (*entity.Todo, error)
	GetTodoByID(ctx context.Context, id, userID uint) (*entity.Todo, error)
	GetUserTodos(ctx context.Context, userID uint, filter entity.TodoFilter) ([]*entity.Todo, int64, error)
	UpdateTodo(ctx context.Context, id uint, title, description string, completed bool, userID uint) (*entity.Todo, error)
	DeleteTodo(ctx context.Context, id, userID uint) error
	CompleteTodo(ctx context.Context, id, userID uint) (*entity.Todo, error)
}

// todoUseCase implements TodoUseCase
type todoUseCase struct {
	todoRepo repository.TodoRepository
}

// NewTodoUseCase creates a new TodoUseCase
func NewTodoUseCase(todoRepo repository.TodoRepository) TodoUseCase {
	return &todoUseCase{
		todoRepo: todoRepo,
	}
}

// CreateTodo creates a new todo
func (uc *todoUseCase) CreateTodo(ctx context.Context, title, description string, userID uint) (*entity.Todo, error) {
	if title == "" {
		return nil, ErrInvalidTodoData
	}

	todo := entity.NewTodo(title, description, userID)
	if err := uc.todoRepo.Create(ctx, todo); err != nil {
		return nil, ErrTodoCreateFailed
	}

	return todo, nil
}

// GetTodoByID retrieves a todo by its ID
func (uc *todoUseCase) GetTodoByID(ctx context.Context, id, userID uint) (*entity.Todo, error) {
	todo, err := uc.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTodoNotFound
	}

	if !todo.BelongsToUser(userID) {
		return nil, ErrNotAuthorized
	}

	return todo, nil
}

// GetUserTodos retrieves todos for a specific user
func (uc *todoUseCase) GetUserTodos(ctx context.Context, userID uint, filter entity.TodoFilter) ([]*entity.Todo, int64, error) {
	filter.UserID = userID
	
	// Ensure pagination defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	todos, err := uc.todoRepo.GetByUserID(ctx, userID, filter)
	if err != nil {
		return nil, 0, err
	}

	count, err := uc.todoRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return todos, count, nil
}

// UpdateTodo updates a todo
func (uc *todoUseCase) UpdateTodo(ctx context.Context, id uint, title, description string, completed bool, userID uint) (*entity.Todo, error) {
	if title == "" {
		return nil, ErrInvalidTodoData
	}

	todo, err := uc.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTodoNotFound
	}

	if !todo.BelongsToUser(userID) {
		return nil, ErrNotAuthorized
	}

	todo.Update(title, description, completed)
	if err := uc.todoRepo.Update(ctx, todo); err != nil {
		return nil, ErrTodoUpdateFailed
	}

	return todo, nil
}

// DeleteTodo deletes a todo
func (uc *todoUseCase) DeleteTodo(ctx context.Context, id, userID uint) error {
	todo, err := uc.todoRepo.GetByID(ctx, id)
	if err != nil {
		return ErrTodoNotFound
	}

	if !todo.BelongsToUser(userID) {
		return ErrNotAuthorized
	}

	if err := uc.todoRepo.Delete(ctx, id); err != nil {
		return ErrTodoDeleteFailed
	}

	return nil
}

// CompleteTodo marks a todo as completed
func (uc *todoUseCase) CompleteTodo(ctx context.Context, id, userID uint) (*entity.Todo, error) {
	todo, err := uc.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTodoNotFound
	}

	if !todo.BelongsToUser(userID) {
		return nil, ErrNotAuthorized
	}

	todo.MarkAsCompleted()
	todo.UpdatedAt = time.Now()

	if err := uc.todoRepo.Update(ctx, todo); err != nil {
		return nil, ErrTodoUpdateFailed
	}

	return todo, nil
}
