package pg_db

import (
	"database/sql"
)

type Store struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB, q *Queries) *Store {
	return &Store{
		db:      db,
		Queries: q,
	}
}

func (store *Store) GetDB() *sql.DB {
	return store.db
}

func (store *Store) Close() error {
	return store.db.Close()
}
