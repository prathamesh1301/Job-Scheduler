package store

import (
	"context"
	"database/sql"
	"time"
)

type RefreshToken struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type RefreshTokenStore interface {
	GetRefreshTokenByToken(ctx context.Context, token string) (*RefreshToken, error)
	InsertRefreshToken(ctx context.Context, refreshToken *RefreshToken) error
	UpdateRefreshToken(ctx context.Context, userID string, token string) error
	EnforceSessionLimit(ctx context.Context, userID string, maxSessions int) error
}

type RefreshTokenStoreImpl struct {
	db *sql.DB
}

func (s *RefreshTokenStoreImpl) GetRefreshTokenByToken(ctx context.Context, token string) (*RefreshToken, error) {
	refreshToken := &RefreshToken{}
	query := `SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token = $1`
	err := s.db.QueryRowContext(ctx, query, token).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func(s *RefreshTokenStoreImpl) InsertRefreshToken(ctx context.Context, refreshToken *RefreshToken) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := s.db.ExecContext(ctx, query, refreshToken.UserID, refreshToken.Token, refreshToken.ExpiresAt)
	if err != nil {
		return err
	}
	return nil
}

func(s *RefreshTokenStoreImpl) UpdateRefreshToken(ctx context.Context, userID string, token string) error {
	query := `UPDATE refresh_tokens SET token = $1 WHERE user_id = $2`
	_, err := s.db.ExecContext(ctx, query, token, userID)
	if err != nil {
		return err
	}
	return nil
}

func(s *RefreshTokenStoreImpl) EnforceSessionLimit(ctx context.Context, userID string, maxSessions int) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE id IN (
			SELECT id FROM refresh_tokens
			WHERE user_id = $1
			ORDER BY id DESC
			OFFSET $2
		)`
	_, err := s.db.ExecContext(ctx, query, userID, maxSessions)
	return err
}