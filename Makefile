#! Database configuration
DB_USER := postgres
DB_PASSWORD := 123321
DB_HOST := localhost
DB_PORT := 6432
DB_TEST_NAME := simple_bank_test
DB_DEV_NAME := simple_bank_dev
DB_URL_PREFIX := postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)

#! Docker configuration
POSTGRES_CONTAINER := postgres
POSTGRES_VERSION := 12-alpine
MIGRATION_DIR := db/migration

#! Help command
.PHONY: help
help:
	@echo "Available commands:"
	@echo "Database Operations:"
	@echo "  start-postgres    - Start PostgreSQL container"
	@echo "  psql             - Connect to PostgreSQL with psql"
	@echo ""
	@echo "Test Database Commands:"
	@echo "  create-test-db   - Create test database"
	@echo "  drop-test-db     - Drop test database"
	@echo "  connect-test-db  - Connect to test database"
	@echo "  migrate-test-up  - Run test database migrations up"
	@echo "  migrate-test-down- Run test database migrations down"
	@echo ""
	@echo "Development Database Commands:"
	@echo "  create-dev-db    - Create development database"
	@echo "  drop-dev-db      - Drop development database"
	@echo "  connect-dev-db   - Connect to development database"
	@echo "  migrate-dev-up   - Run development database migrations up"
	@echo "  migrate-dev-down - Run development database migrations down"
	@echo ""
	@echo "Migration Commands:"
	@echo "  create-migration - Create new migration files"
	@echo ""
	@echo "SQLC Commands:"
	@echo "  pull-sqlc       - Pull SQLC Docker image"
	@echo "  init-sqlc       - Initialize SQLC configuration"
	@echo "  generate-sqlc   - Generate Go code from SQL"
	@echo ""
	@echo "Testing Commands:"
	@echo "  test            - Run tests with coverage"
	@echo "  clean-test      - Clean test cache"

#! PostgreSQL commands
.PHONY: start-postgres psql
start-postgres:
	docker run --name $(POSTGRES_CONTAINER) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-p $(DB_PORT):5432 \
		-d postgres:$(POSTGRES_VERSION)

psql:
	docker exec -it $(POSTGRES_CONTAINER) psql -U $(DB_USER) -d $(DB_USER)

#! Test database commands
.PHONY: create-test-db drop-test-db connect-test-db migrate-test-up migrate-test-down
create-test-db:
	docker exec -it $(POSTGRES_CONTAINER) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_TEST_NAME)

drop-test-db:
	docker exec -it $(POSTGRES_CONTAINER) dropdb -U $(DB_USER) $(DB_TEST_NAME)

connect-test-db:
	docker exec -it $(POSTGRES_CONTAINER) psql --username=$(DB_USER) --dbname=$(DB_USER) $(DB_TEST_NAME)

migrate-test-up:
	migrate -path $(MIGRATION_DIR)/ -database "$(DB_URL_PREFIX)/$(DB_TEST_NAME)?sslmode=disable" --verbose up

migrate-test-down:
	migrate -path $(MIGRATION_DIR)/ -database "$(DB_URL_PREFIX)/$(DB_TEST_NAME)?sslmode=disable" --verbose down

#! Development database commands
.PHONY: create-dev-db drop-dev-db connect-dev-db migrate-dev-up migrate-dev-down
create-dev-db:
	docker exec -it $(POSTGRES_CONTAINER) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_DEV_NAME)

drop-dev-db:
	docker exec -it $(POSTGRES_CONTAINER) dropdb -U $(DB_USER) $(DB_DEV_NAME)

connect-dev-db:
	docker exec -it $(POSTGRES_CONTAINER) psql --username=$(DB_USER) --dbname=$(DB_USER) $(DB_DEV_NAME)

migrate-dev-up:
	migrate -path $(MIGRATION_DIR)/ -database "$(DB_URL_PREFIX)/$(DB_DEV_NAME)?sslmode=disable" --verbose up

migrate-dev-down:
	migrate -path $(MIGRATION_DIR)/ -database "$(DB_URL_PREFIX)/$(DB_DEV_NAME)?sslmode=disable" --verbose down

#! Migration commands
.PHONY: create-migration
create-migration:
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq init_schema

#! SQLC commands
.PHONY: pull-sqlc init-sqlc generate-sqlc
pull-sqlc:
	docker pull sqlc/sqlc

init-sqlc:
	docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc init

generate-sqlc:
	docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc generate

#! Testing commands
.PHONY: test clean-test
test:
	go test -v -cover ./...

clean-test-cache:
	go clean -testcache