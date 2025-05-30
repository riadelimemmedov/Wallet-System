#!/bin/sh

# Exit immediately if a command exits with a non-zero status
set -e

# Function to wait for postgres
wait_for_postgres() {
    local db_host
    local db_port
    
    if [ "$APP_ENV" = "dev" ]; then
        db_host="postgres"
        db_port="$DEV_DB_CONTAINER_PORT"
        DB_SOURCE="$DEV_DB_SOURCE"
    else
        db_host="postgres_test"
        db_port="$TEST_DB_CONTAINER_PORT"
        DB_SOURCE="$TEST_DB_SOURCE"
    fi
    
    echo "Waiting for PostgreSQL ($APP_ENV) at ${db_host}:${db_port}..."
    /usr/local/bin/wait-for.sh "${db_host}:${db_port}" -t 60
    echo "PostgreSQL ($APP_ENV) is up and running!"
}

# Run migrations if MIGRATE_ON_STARTUP is set to true
run_migrations() {
    if [ "$MIGRATE_ON_STARTUP" = "true" ]; then
        echo "Running database migrations for $APP_ENV environment..."
        
        migrate_output=$(migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up 2>&1)
        migrate_status=$?
        
        if echo "$migrate_output" | grep -q "Dirty database version"; then
            echo "Detected dirty database. Attempting to fix..."
            dirty_version=$(echo "$migrate_output" | grep -o "Dirty database version [0-9]*" | awk '{print $4}')
            
            echo "Forcing version $dirty_version to fix dirty state..."
            migrate -path /app/db/migration -database "$DB_SOURCE" force $dirty_version
            
            echo "Retrying migrations..."
            migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up
        elif [ $migrate_status -ne 0 ]; then
            echo "Migration failed with error:"
            echo "$migrate_output"
            exit 1
        fi
        
        echo "Migrations completed!"
    else
        echo "Skipping migrations..."
    fi
}

# Main execution
echo "Starting application in $APP_ENV environment..."

# Wait for PostgreSQL to be ready
wait_for_postgres

# Run migrations if enabled
run_migrations

# Execute the main application
echo "Executing main application..."
exec "$@"