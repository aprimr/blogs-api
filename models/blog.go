package models

import "time"

type Blog struct {
	BlogId      string    `json:"blogid"`
	Uid         string    `json:"uid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	IsDeleted   bool      `json:"is_deleted"`
	IsPrivate   bool      `json:"is_private"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// User sends BlogBody to the server while creating new blog
type BlogBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	IsPrivate   bool   `json:"is_private"`
}
