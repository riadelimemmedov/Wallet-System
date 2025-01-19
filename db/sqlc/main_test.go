package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
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

// !TestConfig holds test configuration options
type TestConfig struct {
	EnvPath string
}

// !getEnvAsInt safely gets on environment variable as integer
func getEnvAsInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return defaultVal
}

// !getEnvAsDuration safely gets on environment variable as duration
func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if v, err := time.ParseDuration(value); err == nil {
			return v
		}
	}
	return defaultVal
}

// ! NewTestDB initializes database connection for testing with custom config
func NewTestDB(config TestConfig) (*TestDB, error) {
	//? Load test environment variables from specified path
	if err := godotenv.Load(config.EnvPath); err != nil {
		return nil, fmt.Errorf("error loading env file from %s: %w", config.EnvPath, err)
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
	dbURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("TEST_DB_SSLMODE"),
	)

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("parsing database URL: %w", err)
	}

	//? Configure connection pool settings from environment variables
	config.MaxConns = int32(getEnvAsInt("TEST_DB_MAX_CONNS", 5))
	config.MinConns = int32(getEnvAsInt("TEST_DB_MIN_CONNS", 2))
	config.MaxConnLifetime = getEnvAsDuration("TEST_DB_CONN_LIFETIME", time.Hour)
	config.MaxConnIdleTime = getEnvAsDuration("TEST_DB_CONN_IDLE_TIME", 30*time.Minute)
	config.HealthCheckPeriod = getEnvAsDuration("TEST_DB_HEALTH_CHECK_PERIOD", 1*time.Minute)

	//? Configure connection attempts
	maxRetries := getEnvAsInt("TEST_DB_MAX_RETRIES", 3)
	var pool *pgxpool.Pool

	for i := 0; i < maxRetries; i++ {
		pool, err = pgxpool.ConnectConfig(ctx, config)
		if err == nil {
			break
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
	}

	//? Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func TestMain(m *testing.M) {
	config := TestConfig{
		EnvPath: "../../.env.test",
	}
	var err error
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
