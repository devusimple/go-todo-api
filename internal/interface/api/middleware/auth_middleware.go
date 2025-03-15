package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"todo-api/internal/util/jwt"
)

// User context key
const (
	UserIDKey     = "user_id"
	UserUsernameKey = "username"
)

// AuthMiddleware is a middleware for authentication
type AuthMiddleware struct {
	jwtService jwt.JWTService
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(jwtService jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// Authenticate authenticates the request
func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Authorization header is required",
			})
		}

		// Check if the token is in the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Authorization header format must be Bearer {token}",
			})
		}

		// Extract token
		tokenString := parts[1]

		// Validate token
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid or expired token",
			})
		}

		// Set user ID in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserUsernameKey, claims.Username)

		// Continue
		return next(c)
	}
}

// GetUserIDFromContext gets the user ID from the context
func GetUserIDFromContext(c echo.Context) uint {
	userID := c.Get(UserIDKey)
	if userID == nil {
		return 0
	}
	return userID.(uint)
}

// GetUsernameFromContext gets the username from the context
func GetUsernameFromContext(c echo.Context) string {
	username := c.Get(UserUsernameKey)
	if username == nil {
		return ""
	}
	return username.(string)
}
