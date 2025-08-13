package domain

import "time"

type Book struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	IsFree      bool      `json:"is_free" binding:"required"`
	Genres      []string  `json:"genres" binding:"required"`
	PublishedAt time.Time `json:"published_at"`
}
