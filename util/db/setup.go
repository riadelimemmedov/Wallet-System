package setup

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	db "github.com/riad/banksystemendtoend/db/sqlc"
	setup "github.com/riad/banksystemendtoend/util/cache"
	"github.com/riad/banksystemendtoend/util/common"
	"github.com/riad/banksystemendtoend/util/config"
	"github.com/riad/banksystemendtoend/util/env"
	"go.uber.org/zap"
)

var (
	store db.Store
	pool  *pgxpool.Pool
	ctx   = context.Background()
	log   = zap.L()
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
		log.Error("Failed to create app environment config", zap.String("environment", string(environment.AppEnv)), zap.Error(err))
		return err
	}

	// Initialize database
	store, err = InitializeDB(ctx, config, environment.Prefix)
	if err != nil {
		log.Error("Failed to initialize database", zap.String("prefix", environment.Prefix), zap.Error(err))
		return nil
	}

	_, err = db.GetSQLStore(store)
	if err != nil {
		log.Error("Failed to get SQL store", zap.Error(err))
		return err
	}

	// Initialize Redis
	_, err = setup.InitializeRedis(config, environment.Prefix)
	if err != nil {
		log.Error("Failed to initialize Redis", zap.String("prefix", environment.Prefix), zap.Error(err))
		return nil
	}

	log.Info("Successfully connected to database", zap.String("environment", string(environment.AppEnv)))
	return nil
}

// GetStore return initialized store for db
func GetStore() db.Store {
	return store
}

// CheckDBHealth checks the health of the database connection
func CheckDBHealth(ctx context.Context, store db.Store) error {
	sqlStore, err := db.GetSQLStore(store)
	if err != nil {
		log.Error("Failed to get SQL store", zap.Error(err))
		return fmt.Errorf("failed to get SQL store: %w", err)
	}

	if sqlStore.Pool == nil {
		log.Error("Database pool is not initialized")
		return fmt.Errorf("database pool is not initialized")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := sqlStore.Pool.Ping(timeoutCtx); err != nil {
		log.Error("Database health check failed", zap.Error(err))
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}
