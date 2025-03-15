package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Custom errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// Claims represents JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTService handles JWT operations
type JWTService interface {
	GenerateToken(userID uint, username string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

// jwtService implements JWTService
type jwtService struct {
	secretKey  string
	expiration time.Duration
}

// NewJWTService creates a new JWTService
func NewJWTService(secretKey string, expiration time.Duration) JWTService {
	return &jwtService{
		secretKey:  secretKey,
		expiration: expiration,
	}
}

// GenerateToken generates a new JWT token
func (s *jwtService) GenerateToken(userID uint, username string) (string, error) {
	// Create claims
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Get claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Check if token is expired
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
