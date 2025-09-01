package domain

import "errors"

var (
	ErrBookNotFound        = errors.New("book not found")
	ErrRefreshTokenExpired = errors.New("session expired")
	ErrUserNotFound        = errors.New("user not found")
)
