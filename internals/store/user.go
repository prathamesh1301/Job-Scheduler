package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Password     []byte `json:"-"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// UserRepository defines the interface for user data access.
// Any implementation (Postgres, mock, etc.) must satisfy this.
type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = $1`
	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserStore) CreateUser(ctx context.Context, user *User) (*User, error) {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, created_at, updated_at`

	err := s.db.QueryRowContext(ctx, query, user.Username, user.Password).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

