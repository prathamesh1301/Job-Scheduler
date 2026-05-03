package store

import (
	"context"
	"database/sql"
)

type IdempotencyStore interface {
	Create(ctx context.Context, idempotencyKey string) error
	Delete(ctx context.Context, idempotencyKey string) error
	Exists(ctx context.Context, idempotencyKey string) (bool, error)
}

type IdempotencyStoreImpl struct {
	db *sql.DB
}

func NewIdempotencyStore(db *sql.DB) *IdempotencyStoreImpl {
	return &IdempotencyStoreImpl{db: db}
}

func (s *IdempotencyStoreImpl) Create(ctx context.Context, idempotencyKey string) error {
	query := `INSERT INTO idempotency (idempotency_key) VALUES ($1)`
	_, err := s.db.ExecContext(ctx, query, idempotencyKey)
	return err
}


func (s *IdempotencyStoreImpl) Delete(ctx context.Context, idempotencyKey string) error {
	query := `DELETE FROM idempotency WHERE idempotency_key = $1`
	_, err := s.db.ExecContext(ctx, query, idempotencyKey)
	return err
}

func (s *IdempotencyStoreImpl) Exists(ctx context.Context, idempotencyKey string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM idempotency WHERE idempotency_key = $1)`
	var exists bool
	err := s.db.QueryRowContext(ctx, query, idempotencyKey).Scan(&exists)
	return exists, err
}