package psql

import (
	"context"
	"database/sql"

	"github.com/jackietana/crud-app/internal/domain"
)

type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepo(db *sql.DB) *TokenRepository {
	return &TokenRepository{db}
}

func (tr *TokenRepository) Create(ctx context.Context, t domain.RefreshToken) error {
	strExec := "INSERT INTO refresh_tokens (user_id, token, expires_at) values ($1, $2, $3)"
	_, err := tr.db.ExecContext(ctx, strExec, t.UserID, t.Token, t.ExpiresAt)

	return err
}

func (tr *TokenRepository) Get(ctx context.Context, token string) (domain.RefreshToken, error) {
	strExec := "SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token=$1"
	var t domain.RefreshToken
	err := tr.db.QueryRowContext(ctx, strExec, token).Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt)
	if err != nil {
		return t, err
	}

	_, err = tr.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE user_id=$1", t.UserID)

	return t, err
}
