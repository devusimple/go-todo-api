package presenter

import (
	"time"

	"todo-api/internal/domain/entity"
)

// UserResponse represents a user response
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserResponse converts a user entity to a user response
func UserResponse(user *entity.User) map[string]interface{} {
	return map[string]interface{}{
		"data": UserResponseData(user),
	}
}

// LoginResponse creates a login response
func LoginResponse(token string, user *entity.User) map[string]interface{} {
	return map[string]interface{}{
		"token": token,
		"user":  UserResponseData(user),
	}
}

// UserResponseData converts a user entity to a user response data
func UserResponseData(user *entity.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
