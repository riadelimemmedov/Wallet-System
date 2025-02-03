package setup

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	db "github.com/riad/banksystemendtoend/db/sqlc"
	"github.com/riad/banksystemendtoend/util/common"
	"github.com/riad/banksystemendtoend/util/config"
	"github.com/riad/banksystemendtoend/util/env"
)

var (
	store db.Store
	ctx   = context.Background()
)

// InitializeDB initializes database connection with custom config
func InitializeDB(ctx context.Context, config config.AppConfig, envPrefix string) (db.Store, error) {
	if config.ConfigFilePath != "" {
		if err := godotenv.Load(config.ConfigFilePath); err != nil {
			return nil, fmt.Errorf("error loading env file from %s: %w", config.ConfigFilePath, err)
		}
	}

	pool, err := setupDBPool(ctx, envPrefix)
	if err != nil {
		return nil, fmt.Errorf("error creating pool: %w", err)
	}

	store := db.NewStore(pool)
	return store, nil
}

// setupDBPool establishes database connection
func setupDBPool(ctx context.Context, envPrefix string) (*pgxpool.Pool, error) {
	return common.SetupDBPool(ctx, envPrefix)
}

// DropAllData removes all data from the database
func DropAllData(ctx context.Context, store db.Store) error {
	sqlStore, err := db.GetSQLStore(store)
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

// InitializeEnvironment environment variable
func InitializeEnvironment(environment config.Environment) error {
	config, err := env.NewAppEnvironmentConfig(environment.AppEnv)
	if err != nil {
		return err
	}

	store, err = InitializeDB(ctx, config, environment.Prefix)
	if err != nil {
		return err
	}

	_, err = db.GetSQLStore(store)
	if err != nil {
		return err
	}

	log.Printf("Successfully connected to %s database", environment)
	return nil
}

func GetStore() db.Store {
	return store
}
