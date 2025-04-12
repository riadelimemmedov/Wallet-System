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
	@echo "Development Commands:"
	@echo "  make docker-dev-build    - Build development containers"
	@echo "  make docker-dev-up       - Start development environment"
	@echo "  make docker-dev-down     - Stop development environment"
	@echo "  make docker-dev-logs     - View development logs"
	@echo ""
	@echo "Test Commands:"
	@echo "  make docker-test-build   - Build test containers"
	@echo "  make docker-test-run     - Run tests and exit"
	@echo "  make docker-test-up      - Start test environment"
	@echo "  make docker-test-down    - Stop test environment"
	@echo "  make docker-test-logs    - View test logs"
	@echo ""
	@echo "Database Commands:"
	@echo "  make docker-db-dev       - Connect to development database"
	@echo "  make docker-db-test      - Connect to test database"
	@echo "  make docker-migrate-up   - Run database migrations"
	@echo "  make docker-migrate-down - Rollback database migrations"
	@echo ""
	@echo "Cleanup Commands:"
	@echo "  make docker-clean        - Remove all containers and volumes"
	@echo "  make docker-clean-all    - Remove everything including images"
	@echo "  make docker-prune        - Remove unused images and volumes"
	@echo ""
	@echo "Utility Commands:"
	@echo "  make docker-ps           - Show running containers"
	@echo "  make docker-logs         - View container logs (CONTAINER=name)"

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

#! Testing purpose command for db migrations
init-sqlc:
	docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc init

generate-sqlc:
	docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc generate

create_migration:
	migrate create -ext sql -dir db/migration -seq init_schema
migrate_up:
	migrate -path db/migration/ -database "postgresql://postgres:123321@localhost:6432/simple_bank?sslmode=disable" --verbose up 2(migration folder example)
	migrate -path db/migration/ -database "postgresql://postgres:123321@localhost:6432/simple_bank_test?sslmode=disable" force 1
migrate_down:
	migrate -path db/migration/ -database "postgresql://postgres:123321@localhost:6432/simple_bank?sslmode=disable" --verbose down

#! Testing commands
.PHONY: test clean-test
test:
	go test -v -cover ./...

clean-test-cache:
	go clean -testcache


#! Api
.PHONY: run-api
run-api:
	air


#! Docker Development Commands
.PHONY: docker-dev-build docker-dev-up docker-dev-down docker-dev-logs

#? Build development containers
docker-dev-build:
	docker-compose -f docker/development/docker-compose.yml \
		--env-file docker/development/.env build

#? Start development environment
docker-dev-up:
	docker-compose -f docker/development/docker-compose.yml \
		--env-file docker/development/.env \
		up -d

#? Stop development environment and remove volumes
docker-dev-down:
	docker-compose -f docker/development/docker-compose.yml \
		--env-file docker/development/.env down -v

#? View development logs
docker-dev-logs:
	docker-compose -f docker/development/docker-compose.yml \
		--env-file docker/development/.env up \
		logs -f

#! Docker Test Commands
.PHONY: docker-test-build docker-test-up docker-test-down docker-test-logs docker-test-run

#? Build test containers
docker-test-build:
	docker-compose -f docker/test/docker-compose.yml \
		--env-file docker/development/.env.test up \
		build

#? Start test environment
docker-test-up:
	docker-compose -f docker/test/docker-compose.yml \
		--env-file docker/development/.env.test up \
		up -d

#? Stop test environment and remove volumes
docker-test-down:
	docker-compose -f docker/test/docker-compose.yml \
		--env-file docker/development/.env.test up \
		down -v

#? View test logs
docker-test-logs:
	docker-compose -f docker/test/docker-compose.yml \
		--env-file docker/development/.env.test up \
		logs -f

#? Run tests and exit
docker-test-run:
	docker-compose -f docker/test/docker-compose.yml \
		--env-file .env.test \
		up --abort-on-container-exit --exit-code-from api_test

#! Docker Database Commands
.PHONY: docker-db-dev docker-db-test

DEV_DB_USER :=postgres
DEV_DB_NAME :=simple_bank_dev
DEV_DB_CONTAINER_PORT := 5433
TEST_DB_USER :=postgres
TEST_DB_NAME :=simple_bank_test
TEST_DB_CONTAINER_PORT := 5434

#? Connect to development database
docker-db-dev:
	docker exec -it simple_bank_dev_db psql -U ${DEV_DB_USER} -d ${DEV_DB_NAME} -p ${DEV_DB_CONTAINER_PORT}

#? Connect to test database
docker-db-test:
	docker exec -it simple_bank_test_db psql -U ${TEST_DB_USER} -d ${TEST_DB_NAME} -p ${TEST_DB_CONTAINER_PORT}

#! Docker Clean Commands
.PHONY: docker-clean docker-clean-all docker-prune

#? Stop and remove all containers and volumes
docker-clean:
	docker-compose -f docker/development/docker-compose.yml down -v
	docker-compose -f docker/test/docker-compose.yml down -v

#? Remove all unused containers, networks, images and volumes
docker-clean-all: docker-clean
	docker system prune -af --volumes

#? Remove only dangling images and unused volumes
docker-prune:
	docker image prune -f
	docker volume prune -f

#! Docker Utility Commands
.PHONY: docker-ps docker-logs

#? Show running containers
docker-ps:
	docker ps

#? View logs for a specific container (usage: make docker-logs CONTAINER=container_name)
docker-logs:
	@if [ "$(CONTAINER)" = "" ]; then \
		echo "Please specify a container name: make docker-logs CONTAINER=container_name"; \
		exit 1; \
	fi
	docker logs -f $(CONTAINER)


# docker exec simple_bank_dev_api migrate -path db/migration/ -database "postgresql://postgres:123321@postgres:5433/simple_bank_dev?sslmode=disable" --verbose up
# docker exec -it b62 redis-cli
# docker ps | findstr redis

# docker exec -it b62 redis-cli -a 6V@b)B5V)V2K