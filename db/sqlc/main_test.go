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

// type TestDB struct {
// 	Pool    *pgxpool.Pool
// 	Queries *Queries
// 	cleanup func()
// }

var (
	testStore Store
	ctx       = context.Background()
)

const appEnvironment = "test"
const envPrefix = "TEST"

// ! NewTestDB initializes database connection for testing with custom config
func NewTestDB(config config.AppConfig) (Store, error) {
	//? Load test environment variables from specified path
	if err := godotenv.Load(config.ConfigFilePath); err != nil {
		return nil, fmt.Errorf("error loading env file from %s: %w", config.ConfigFilePath, err)
	}

	pool, err := setupTestPool()
	if err != nil {
		return nil, fmt.Errorf("error creating test pool: %w", err)
	}

	testStore := NewStore(pool)
	return testStore, nil
}

// !setupTestPool check db connection
func setupTestPool() (*pgxpool.Pool, error) {
	return common.SetupDBPool(ctx, envPrefix)
}

// !DropAllData removes all data from the database
func DropAllData() error {
	sqlStore, err := GetSQLStore(testStore)
	if err != nil {
		return err
	}
	query := common.GetTruncateTablesQuery()
	_, err = sqlStore.Pool.Exec(ctx, query)
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
	testStore, err = NewTestDB(config)
	sqlStore, err := GetSQLStore(testStore)
	if err != nil {
		log.Fatalf("Failed to get SQL store: %v", err)
	}
	fmt.Println("Tested connection to db...")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	code := m.Run()
	if sqlStore != nil && sqlStore.cleanup != nil {
		sqlStore.cleanup()
	}
	os.Exit(code)
}

// ! CleanupDB provides a helper function for cleaning up the test database betweesn tests.
func CleanupDB(t *testing.T) {
	t.Helper()
	if err := DropAllData(); err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}
}
