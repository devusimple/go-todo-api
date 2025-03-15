package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"todo-api/internal/domain/entity"
	"todo-api/internal/domain/usecase"
	"todo-api/internal/interface/api/middleware"
	"todo-api/internal/interface/api/presenter"
)

// TodoHandler handles HTTP requests related to todos
type TodoHandler struct {
	todoUseCase usecase.TodoUseCase
	logger      *logrus.Logger
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(todoUseCase usecase.TodoUseCase, logger *logrus.Logger) *TodoHandler {
	return &TodoHandler{
		todoUseCase: todoUseCase,
		logger:      logger,
	}
}

// CreateTodoRequest represents the request to create a todo
type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description"`
}

// UpdateTodoRequest represents the request to update a todo
type UpdateTodoRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// CreateTodo handles the creation of a new todo
func (h *TodoHandler) CreateTodo(c echo.Context) error {
	// Get user ID from context
	userID := middleware.GetUserIDFromContext(c)

	// Parse request
	req := new(CreateTodoRequest)
	if err := c.Bind(req); err != nil {
		h.logger.WithError(err).Error("Failed to bind request")
		return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid request"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		h.logger.WithError(err).Error("Invalid request payload")
		return c.JSON(http.StatusBadRequest, presenter.ValidationErrorResponse(err))
	}

	// Create todo
	todo, err := h.todoUseCase.CreateTodo(c.Request().Context(), req.Title, req.Description, userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create todo")
		return c.JSON(http.StatusInternalServerError, presenter.ErrorResponse("Failed to create todo"))
	}

	// Return response
	return c.JSON(http.StatusCreated, presenter.TodoResponse(todo))
}

// GetTodo handles retrieving a todo by ID
func (h *TodoHandler) GetTodo(c echo.Context) error {
	// Get user ID from context
	userID := middleware.GetUserIDFromContext(c)

	// Parse todo ID
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.WithError(err).Error("Invalid todo ID")
		return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid todo ID"))
	}

	// Get todo
	todo, err := h.todoUseCase.GetTodoByID(c.Request().Context(), uint(todoID), userID)
	if err != nil {
		switch err {
		case usecase.ErrTodoNotFound:
			return c.JSON(http.StatusNotFound, presenter.ErrorResponse("Todo not found"))
		case usecase.ErrNotAuthorized:
			return c.JSON(http.StatusForbidden, presenter.ErrorResponse("Not authorized to access this todo"))
		default:
			h.logger.WithError(err).Error("Failed to get todo")
			return c.JSON(http.StatusInternalServerError, presenter.ErrorResponse("Failed to get todo"))
		}
	}

	// Return response
	return c.JSON(http.StatusOK, presenter.TodoResponse(todo))
}

// GetTodos handles retrieving all todos for a user
func (h *TodoHandler) GetTodos(c echo.Context) error {
	// Get user ID from context
	userID := middleware.GetUserIDFromContext(c)

	// Parse filter parameters
	filter := entity.TodoFilter{}

	// Parse completed filter
	completed := c.QueryParam("completed")
	if completed != "" {
		completedBool, err := strconv.ParseBool(completed)
		if err == nil {
			filter.Completed = &completedBool
		}
	}

	// Parse search query
	filter.Search = c.QueryParam("search")

	// Parse pagination
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	filter.Page = page

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	filter.PageSize = pageSize

	// Get todos
	todos, count, err := h.todoUseCase.GetUserTodos(c.Request().Context(), userID, filter)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get todos")
		return c.JSON(http.StatusInternalServerError, presenter.ErrorResponse("Failed to get todos"))
	}

	// Return response
	return c.JSON(http.StatusOK, presenter.TodosResponse(todos, count, page, pageSize))
}

// UpdateTodo handles updating a todo
func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	// Get user ID from context
	userID := middleware.GetUserIDFromContext(c)

	// Parse todo ID
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.WithError(err).Error("Invalid todo ID")
		return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid todo ID"))
	}

	// Parse request
	req := new(UpdateTodoRequest)
	if err := c.Bind(req); err != nil {
		h.logger.WithError(err).Error("Failed to bind request")
		return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid request"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		h.logger.WithError(err).Error("Invalid request payload")
		return c.JSON(http.StatusBadRequest, presenter.ValidationErrorResponse(err))
	}

	// Update todo
	todo, err := h.todoUseCase.UpdateTodo(c.Request().Context(), uint(todoID), req.Title, req.Description, req.Completed, userID)
	if err != nil {
		switch err {
		case usecase.ErrTodoNotFound:
			return c.JSON(http.StatusNotFound, presenter.ErrorResponse("Todo not found"))
		case usecase.ErrNotAuthorized:
			return c.JSON(http.StatusForbidden, presenter.ErrorResponse("Not authorized to update this todo"))
		case usecase.ErrInvalidTodoData:
			return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid todo data"))
		default:
			h.logger.WithError(err).Error("Failed to update todo")
			return c.JSON(http.StatusInternalServerError, presenter.ErrorResponse("Failed to update todo"))
		}
	}

	// Return response
	return c.JSON(http.StatusOK, presenter.TodoResponse(todo))
}

// DeleteTodo handles deleting a todo
func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	// Get user ID from context
	userID := middleware.GetUserIDFromContext(c)

	// Parse todo ID
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.WithError(err).Error("Invalid todo ID")
		return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid todo ID"))
	}

	// Delete todo
	err = h.todoUseCase.DeleteTodo(c.Request().Context(), uint(todoID), userID)
	if err != nil {
		switch err {
		case usecase.ErrTodoNotFound:
			return c.JSON(http.StatusNotFound, presenter.ErrorResponse("Todo not found"))
		case usecase.ErrNotAuthorized:
			return c.JSON(http.StatusForbidden, presenter.ErrorResponse("Not authorized to delete this todo"))
		default:
			h.logger.WithError(err).Error("Failed to delete todo")
			return c.JSON(http.StatusInternalServerError, presenter.ErrorResponse("Failed to delete todo"))
		}
	}

	// Return response
	return c.NoContent(http.StatusNoContent)
}

// CompleteTodo handles marking a todo as completed
func (h *TodoHandler) CompleteTodo(c echo.Context) error {
	// Get user ID from context
	userID := middleware.GetUserIDFromContext(c)

	// Parse todo ID
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.WithError(err).Error("Invalid todo ID")
		return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid todo ID"))
	}

	// Complete todo
	todo, err := h.todoUseCase.CompleteTodo(c.Request().Context(), uint(todoID), userID)
	if err != nil {
		switch err {
		case usecase.ErrTodoNotFound:
			return c.JSON(http.StatusNotFound, presenter.ErrorResponse("Todo not found"))
		case usecase.ErrNotAuthorized:
			return c.JSON(http.StatusForbidden, presenter.ErrorResponse("Not authorized to update this todo"))
		default:
			h.logger.WithError(err).Error("Failed to complete todo")
			return c.JSON(http.StatusInternalServerError, presenter.ErrorResponse("Failed to complete todo"))
		}
	}

	// Return response
	return c.JSON(http.StatusOK, presenter.TodoResponse(todo))
}
