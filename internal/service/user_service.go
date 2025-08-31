package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackietana/crud-app/internal/domain"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
}

type UserService struct {
	repo   UserRepository
	hasher PasswordHasher

	hmacSecret []byte
	tokenTTL   time.Duration
}

func NewUserService(repo UserRepository, hasher PasswordHasher, secret []byte, ttl time.Duration) *UserService {
	return &UserService{repo, hasher, secret, ttl}
}

func (us *UserService) SignUp(ctx context.Context, input domain.User) error {
	password, err := us.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return us.repo.CreateUser(ctx, user)
}

func (us *UserService) SignIn(ctx context.Context, input domain.UserSignIn) (string, error) {
	password, err := us.hasher.Hash(input.Password)
	if err != nil {
		return "", err
	}

	user, err := us.repo.GetByCredentials(ctx, input.Email, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrUserNotFound
		}

		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(user.ID),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(us.tokenTTL).Unix(),
	})

	return token.SignedString(us.hmacSecret)
}

func (us *UserService) ParseToken(ctx context.Context, token string) (int, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return us.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return id, nil
}
