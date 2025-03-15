package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"todo-api/internal/domain/usecase"
	"todo-api/internal/interface/api/middleware"
	"todo-api/internal/interface/api/presenter"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userUseCase usecase.UserUseCase
	logger      *logrus.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userUseCase usecase.UserUseCase, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

// RegisterRequest represents the request to register a new user
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

// LoginRequest represents the request to login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateProfileRequest represents the request to update a user's profile
type UpdateProfileRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
}

// UpdatePasswordRequest represents the request to update a user's password
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=100"`
}

// Register handles user registration
func (h *UserHandler) Register(c echo.Context) error {
	// Parse request
	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		h.logger.WithError(err).Error("Failed to bind request")
		return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid request"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		h.logger.WithError(err).Error("Invalid request payload")
		return c.JSON(http.StatusBadRequest, presenter.ValidationErrorResponse(err))
	}

	// Register user
	user, err := h.userUseCase.Register(c.Request().Context(), req.Username, req.Email, req.Password)
	if err != nil {
		switch err {
		case usecase.ErrUsernameExists:
			return c.JSON(http.StatusConflict, presenter.ErrorResponse("Username already exists"))
		case usecase.ErrEmailExists:
			return c.JSON(http.StatusConflict, presenter.ErrorResponse("Email already exists"))
		case usecase.ErrInvalidUserData:
			return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid user data"))
		default:
			h.logger.WithError(err).Error("Failed to register user")
			return c.JSON(http.StatusInternalServerError, presenter.ErrorResponse("Failed to register user"))
		}
	}

	// Return response
	return c.JSON(http.StatusCreated, presenter.UserResponse(user))
}

// Login handles user login
func (h *UserHandler) Login(c echo.Context) error {
	// Parse request
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		h.logger.WithError(err).Error("Failed to bind request")
		return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid request"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		h.logger.WithError(err).Error("Invalid request payload")
		return c.JSON(http.StatusBadRequest, presenter.ValidationErrorResponse(err))
	}

	// Login
	response, err := h.userUseCase.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case usecase.ErrInvalidCredentials:
			return c.JSON(http.StatusUnauthorized, presenter.ErrorResponse("Invalid credentials"))
		default:
			h.logger.WithError(err).Error("Failed to login")
			return c.JSON(http.StatusInternalServerError, presenter.ErrorResponse("Failed to login"))
		}
	}

	// Return response
	return c.JSON(http.StatusOK, presenter.LoginResponse(response.Token, response.User))
}

// GetProfile handles retrieving the current user's profile
func (h *UserHandler) GetProfile(c echo.Context) error {
	// Get user ID from context
	userID := middleware.GetUserIDFromContext(c)

	// Get user
	user, err := h.userUseCase.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			return c.JSON(http.StatusNotFound, presenter.ErrorResponse("User not found"))
		default:
			h.logger.WithError(err).Error("Failed to get user profile")
			return c.JSON(http.StatusInternalServerError, presenter.ErrorResponse("Failed to get user profile"))
		}
	}

	// Return response
	return c.JSON(http.StatusOK, presenter.UserResponse(user))
}

// UpdateProfile handles updating the current user's profile
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	// Get user ID from context
	userID := middleware.GetUserIDFromContext(c)

	// Parse request
	req := new(UpdateProfileRequest)
	if err := c.Bind(req); err != nil {
		h.logger.WithError(err).Error("Failed to bind request")
		return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid request"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		h.logger.WithError(err).Error("Invalid request payload")
		return c.JSON(http.StatusBadRequest, presenter.ValidationErrorResponse(err))
	}

	// Update profile
	user, err := h.userUseCase.UpdateProfile(c.Request().Context(), userID, req.Username, req.Email)
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			return c.JSON(http.StatusNotFound, presenter.ErrorResponse("User not found"))
		case usecase.ErrUsernameExists:
			return c.JSON(http.StatusConflict, presenter.ErrorResponse("Username already exists"))
		case usecase.ErrEmailExists:
			return c.JSON(http.StatusConflict, presenter.ErrorResponse("Email already exists"))
		case usecase.ErrInvalidUserData:
			return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid user data"))
		default:
			h.logger.WithError(err).Error("Failed to update profile")
			return c.JSON(http.StatusInternalServerError, presenter.ErrorResponse("Failed to update profile"))
		}
	}

	// Return response
	return c.JSON(http.StatusOK, presenter.UserResponse(user))
}

// UpdatePassword handles updating the current user's password
func (h *UserHandler) UpdatePassword(c echo.Context) error {
	// Get user ID from context
	userID := middleware.GetUserIDFromContext(c)

	// Parse request
	req := new(UpdatePasswordRequest)
	if err := c.Bind(req); err != nil {
		h.logger.WithError(err).Error("Failed to bind request")
		return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid request"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		h.logger.WithError(err).Error("Invalid request payload")
		return c.JSON(http.StatusBadRequest, presenter.ValidationErrorResponse(err))
	}

	// Update password
	err := h.userUseCase.UpdatePassword(c.Request().Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			return c.JSON(http.StatusNotFound, presenter.ErrorResponse("User not found"))
		case usecase.ErrInvalidCredentials:
			return c.JSON(http.StatusUnauthorized, presenter.ErrorResponse("Current password is incorrect"))
		case usecase.ErrInvalidUserData:
			return c.JSON(http.StatusBadRequest, presenter.ErrorResponse("Invalid password data"))
		default:
			h.logger.WithError(err).Error("Failed to update password")
			return c.JSON(http.StatusInternalServerError, presenter.ErrorResponse("Failed to update password"))
		}
	}

	// Return response
	return c.NoContent(http.StatusNoContent)
}
