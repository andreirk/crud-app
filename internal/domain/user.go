package domain

import (
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type UserSignIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=5"`
}

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name" validate:"required,gte=2"`
	Email        string    `json:"email" validate:"required,email"`
	Password     string    `json:"password" validate:"required,gte=5"`
	RegisteredAt time.Time `json:"registered_at"`
}
