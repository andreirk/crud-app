package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackietana/crud-app/internal/domain"
	log "github.com/sirupsen/logrus"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user domain.User) error {
	strExec := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"
	_, err := ur.db.ExecContext(ctx, strExec, user.Name, user.Email, user.Password)

	log.Info("Repository: CreateUser")

	return err
}

func (ur *UserRepository) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var u domain.User
	err := ur.db.QueryRowContext(ctx, "SELECT id, name, email, registered_at FROM users WHERE email=$1 AND password=$2", email, password).
		Scan(&u.ID, &u.Name, &u.Email, &u.RegisteredAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, domain.ErrUserNotFound
		}

		return u, err
	}

	log.WithField("id", u.ID).Info("Repository: GetByCredentials")

	return u, err
}
