package entity

import (
	"time"
)

// Todo represents a todo item entity
type Todo struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	Completed   bool      `gorm:"default:false"`
	UserID      uint      `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserID"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// TodoFilter represents the filters for querying todos
type TodoFilter struct {
	UserID    uint
	Completed *bool
	Search    string
	Page      int
	PageSize  int
}

// NewTodo creates a new Todo entity
func NewTodo(title, description string, userID uint) *Todo {
	return &Todo{
		Title:       title,
		Description: description,
		UserID:      userID,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// MarkAsCompleted marks the todo as completed
func (t *Todo) MarkAsCompleted() {
	t.Completed = true
	t.UpdatedAt = time.Now()
}

// MarkAsIncomplete marks the todo as incomplete
func (t *Todo) MarkAsIncomplete() {
	t.Completed = false
	t.UpdatedAt = time.Now()
}

// Update updates the todo with the provided details
func (t *Todo) Update(title, description string, completed bool) {
	t.Title = title
	t.Description = description
	t.Completed = completed
	t.UpdatedAt = time.Now()
}

// BelongsToUser checks if the todo belongs to the specified user
func (t *Todo) BelongsToUser(userID uint) bool {
	return t.UserID == userID
}
