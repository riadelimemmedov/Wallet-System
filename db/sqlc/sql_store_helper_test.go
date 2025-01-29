package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ! SetupTestStore initializes and returns a SQLStore instance for testing.
func SetupTestStore(t *testing.T) *SQLStore {
	sqlStore, err := GetSQLStore(testStore)
	require.NoError(t, err)
	return sqlStore
}
