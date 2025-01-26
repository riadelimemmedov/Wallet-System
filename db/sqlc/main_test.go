package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/riad/banksystemendtoend/util/common"
	"github.com/riad/banksystemendtoend/util/config"
	"github.com/riad/banksystemendtoend/util/env"
)

type TestDB struct {
	Pool    *pgxpool.Pool
	Queries *Queries
	cleanup func()
}

var (
	testDB *TestDB
	ctx    = context.Background()
)

const appEnvironment = "test"
const envPrefix = "TEST"

// ! NewTestDB initializes database connection for testing with custom config
func NewTestDB(config config.AppConfig) (*TestDB, error) {
	//? Load test environment variables from specified path
	if err := godotenv.Load(config.ConfigFilePath); err != nil {
		return nil, fmt.Errorf("error loading env file from %s: %w", config.ConfigFilePath, err)
	}

	pool, err := setupTestPool()
	if err != nil {
		return nil, fmt.Errorf("error creating test pool: %w", err)
	}

	db := &TestDB{
		Pool:    pool,
		Queries: New(pool),
		cleanup: func() { pool.Close() },
	}
	return db, nil
}

// !setupTestPool check db connection
func setupTestPool() (*pgxpool.Pool, error) {
	return common.SetupDBPool(ctx, envPrefix)
}

// !DropAllData removes all data from the database
func (db *TestDB) DropAllData() error {
	query := common.GetTruncateTablesQuery()
	_, err := db.Pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to truncate tables: %w", err)
	}
	return nil
}

// !TestMain serves as the entry point for testing, managing database initialization and cleanup.
func TestMain(m *testing.M) {
	config, err := env.NewAppEnvironmentConfig(appEnvironment)
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	testDB, err = NewTestDB(config)
	fmt.Println("Tested connection to db...")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	code := m.Run()
	if testDB != nil && testDB.cleanup != nil {
		testDB.cleanup()
	}
	os.Exit(code)
}

// ! CleanupDB provides a helper function for cleaning up the test database between tests.
func CleanupDB(t *testing.T, db *TestDB) {
	t.Helper()
	if err := db.DropAllData(); err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}
}
