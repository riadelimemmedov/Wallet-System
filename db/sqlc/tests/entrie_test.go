package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/util/common"
	"github.com/stretchr/testify/require"
)

// Helper function to create a random entry
func createRandomEntry(t *testing.T) db.Entry {
	sqlStore := SetupTestStore(t)

	amount, err := common.RandomNumeric()
	require.NoError(t, err)

	arg := db.CreateEntryParams{
		AccountID: sql.NullInt32{Int32: createRandomAccount(t).AccountID, Valid: true},
		Amount:    amount,
	}

	entry, err := sqlStore.Queries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
	defer CleanupDB(t)
}

func TestGetEntry(t *testing.T) {
	sqlStore := SetupTestStore(t)

	entry1 := createRandomEntry(t)
	require.NotEmpty(t, entry1)

	entry2, err := sqlStore.Queries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

	defer CleanupDB(t)
}
