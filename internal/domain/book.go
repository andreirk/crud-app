package domain

import "time"

type Book struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	IsFree      bool      `json:"is_free"`
	Genres      []string  `json:"genres"`
	PublishedAt time.Time `json:"published_at"`
}

type BookCreate struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	IsFree      *bool    `json:"is_free"`
	Genres      []string `json:"genres"`
}
