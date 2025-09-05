package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackietana/crud-app/internal/domain"
	logger "github.com/jackietana/grpc-logger/pkg/domain"
	log "github.com/sirupsen/logrus"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
}

type TokenRepository interface {
	Create(ctx context.Context, token domain.RefreshToken) error
	Get(ctx context.Context, token string) (domain.RefreshToken, error)
}

type LoggerClient interface {
	SendLogRequest(ctx context.Context, req logger.LogItem) error
}

type UserService struct {
	userRepo     UserRepository
	tokenRepo    TokenRepository
	hasher       PasswordHasher
	loggerClient LoggerClient

	hmacSecret []byte
	tokenTTL   time.Duration
	refreshTTL time.Duration
}

func NewUserService(userRepo UserRepository, tokenRepo TokenRepository, hasher PasswordHasher,
	logger LoggerClient, secret []byte, tokenTTL time.Duration, refreshTTL time.Duration) *UserService {
	return &UserService{userRepo, tokenRepo, hasher, logger, secret, tokenTTL, refreshTTL}
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

	if err := us.userRepo.CreateUser(ctx, user); err != nil {
		return err
	}

	user, err = us.userRepo.GetByCredentials(ctx, input.Email, password)
	if err != nil {
		return err
	}

	if err := us.loggerClient.SendLogRequest(ctx, logger.LogItem{
		Action:    logger.ACTION_REGISTER,
		Entity:    logger.ENTITY_USER,
		EntityID:  int64(user.ID),
		Timestamp: time.Now(),
	}); err != nil {
		log.WithField("service", "User.signUp").Error(err)
	}

	return nil
}

func (us *UserService) SignIn(ctx context.Context, input domain.UserSignIn) (string, string, error) {
	password, err := us.hasher.Hash(input.Password)
	if err != nil {
		return "", "", err
	}

	user, err := us.userRepo.GetByCredentials(ctx, input.Email, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", domain.ErrUserNotFound
		}

		return "", "", err
	}

	return us.generateTokens(ctx, user.ID)
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

func (us *UserService) RefreshTokens(ctx context.Context, strRefreshToken string) (string, string, error) {
	refreshToken, err := us.tokenRepo.Get(ctx, strRefreshToken)
	if err != nil {
		return "", "", err
	}

	if refreshToken.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", domain.ErrRefreshTokenExpired
	}

	return us.generateTokens(ctx, refreshToken.UserID)
}

func (us *UserService) generateTokens(ctx context.Context, userId int) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(userId)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(us.tokenTTL).Unix(),
	})

	accessToken, err := t.SignedString(us.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := us.tokenRepo.Create(ctx, domain.RefreshToken{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(us.refreshTTL),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
