package store

import "database/sql"

type Store struct {
	User UserRepository
	RefreshToken RefreshTokenStore
	Idempotency IdempotencyStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		User: &UserStore{db: db},
		RefreshToken: &RefreshTokenStoreImpl{db: db},
		Idempotency: &IdempotencyStoreImpl{db: db},
	}
}