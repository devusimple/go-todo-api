package presenter

import (
	"time"

	"todo-api/internal/domain/entity"
)

// TodoResponse represents a todo response
type TodoResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TodosListResponse represents a paginated list of todos
type TodosListResponse struct {
	Data       []TodoResponse `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationMeta contains metadata about pagination
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// TodoResponse converts a todo entity to a todo response
func TodoResponse(todo *entity.Todo) map[string]interface{} {
	return map[string]interface{}{
		"data": TodoResponseData(todo),
	}
}

// TodosResponse converts a list of todo entities to a todos response
func TodosResponse(todos []*entity.Todo, totalCount int64, currentPage, pageSize int) map[string]interface{} {
	var todoResponses []TodoResponse
	for _, todo := range todos {
		todoResponses = append(todoResponses, TodoResponseData(todo))
	}

	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}

	return map[string]interface{}{
		"data": todoResponses,
		"pagination": PaginationMeta{
			CurrentPage: currentPage,
			PageSize:    pageSize,
			TotalItems:  totalCount,
			TotalPages:  totalPages,
		},
	}
}

// TodoResponseData converts a todo entity to a todo response data
func TodoResponseData(todo *entity.Todo) TodoResponse {
	return TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		UserID:      todo.UserID,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}
