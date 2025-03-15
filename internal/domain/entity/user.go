package entity

import (
	"time"
)

// User represents a user entity
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"size:100;uniqueIndex;not null"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	Password  string    `gorm:"size:255;not null"`
	Todos     []Todo    `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// NewUser creates a new User entity
func NewUser(username, email, password string) *User {
	return &User{
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(username, email string) {
	u.Username = username
	u.Email = email
	u.UpdatedAt = time.Now()
}

// UpdatePassword updates the user's password
func (u *User) UpdatePassword(hashedPassword string) {
	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
}
