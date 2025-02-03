package db

import (
	"testing"

	db "github.com/riad/banksystemendtoend/db/sqlc"
	setup "github.com/riad/banksystemendtoend/util/db"
	"github.com/stretchr/testify/require"
)

// ! SetupTestStore initializes and returns a SQLStore instance for testing.
func SetupTestStore(t *testing.T) *db.SQLStore {
	sqlStore, err := db.GetSQLStore(setup.GetStore())
	require.NotEmpty(t, sqlStore)
	require.NoError(t, err)
	return sqlStore
}
