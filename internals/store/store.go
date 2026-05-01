package store

import "database/sql"

type Store struct {
	User UserRepository
	RefreshToken RefreshTokenStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		User: &UserStore{db: db},
		RefreshToken: &RefreshTokenStoreImpl{db: db},
	}
}