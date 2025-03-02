package common

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// ! SetupDBPool establishes a connection pool to PostgreSQL using environment-based configuration.
func SetupDBPool(ctx context.Context, envPrefix string) (*pgxpool.Pool, error) {
	dbURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv(envPrefix+"_DB_USER"),
		os.Getenv(envPrefix+"_DB_PASSWORD"),
		os.Getenv(envPrefix+"_DB_HOST"),
		os.Getenv(envPrefix+"_DB_CONTAINER_PORT"),
		os.Getenv(envPrefix+"_DB_NAME"),
		os.Getenv(envPrefix+"_DB_SSLMODE"),
	)

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("parsing database URL: %w", err)
	}

	config.MaxConns = int32(GetEnvAsInt(envPrefix+"_DB_MAX_CONNS", 5))
	config.MinConns = int32(GetEnvAsInt(envPrefix+"_DB_MIN_CONNS", 2))
	config.MaxConnLifetime = GetEnvAsDuration(envPrefix+"_DB_CONN_LIFETIME", time.Hour)
	config.MaxConnIdleTime = GetEnvAsDuration(envPrefix+"_DB_CONN_IDLE_TIME", 30*time.Minute)
	config.HealthCheckPeriod = GetEnvAsDuration(envPrefix+"_DB_HEALTH_CHECK_PERIOD", 1*time.Minute)

	maxRetries := GetEnvAsInt(envPrefix+"_DB_MAX_RETRIES", 3)
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

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

// ! GetTruncateTablesQuery return truncate tables query
func GetTruncateTablesQuery() string {
	return `DO $$
    DECLARE
        r RECORD;
    BEGIN
        FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
            EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
        END LOOP;
    END $$;`
}
