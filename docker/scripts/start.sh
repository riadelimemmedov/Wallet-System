#!/bin/sh

#! Exit immediately if a command exits with a non-zero status
set -e

#! Function to wait for postgres
wait_for_postgres() {
    echo "Waiting for PostgreSQL..."
    /usr/local/bin/wait-for.sh postgres:5432 -t 60
    echo "PostgreSQL is up and running!"
}

#! Run migrations if MIGRATE_ON_STARTUP is set to true
run_migrations() {
    if [ "$MIGRATE_ON_STARTUP" = "true" ]; then
        echo "Running database migrations..."
        migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up
        echo "Migrations completed!"
    else
        echo "Skipping migrations..."
    fi
}

#! Main execution
echo "Starting application..."

#! Wait for PostgreSQL to be ready
wait_for_postgres

#! Run migrations if enabled
run_migrations

#! Execute the main application
echo "Executing main application..."
exec "$@"