package models

import "time"

type Comment struct {
	ID        int       `json:"id"`                              // Unique identifier for the comment
	Text      string    `json:"text" binding:"required,max=150"` // The text of the comment (max 150 characters)
	CreatedAt time.Time `json:"created_at"`
}
