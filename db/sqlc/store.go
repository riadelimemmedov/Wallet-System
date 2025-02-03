package db

import (
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// !Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
}

// ! SqlStore provides all functions to execute db queries and transactions
type SQLStore struct {
	Pool *pgxpool.Pool
	*Queries
	cleanup func()
}

// !NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		Pool:    connPool,
		Queries: New(connPool),
		cleanup: func() { connPool.Close() },
	}
}

// ! Method to get SQLStore safely
func GetSQLStore(store Store) (*SQLStore, error) {
	sqlStore, ok := store.(*SQLStore)
	if !ok {
		return nil, fmt.Errorf("invalid store type")
	}
	return sqlStore, nil
}
