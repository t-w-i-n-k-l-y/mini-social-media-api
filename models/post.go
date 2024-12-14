package models

import "time"

// Post represents a social media post with content(text), likes, and associated comments
type Post struct {
	ID        int       `json:"id"` // Unique identifier for the post
	Content   string    `json:"content" binding:"required,max=250"`
	Likes     int       `json:"likes"`
	Comments  []Comment `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
