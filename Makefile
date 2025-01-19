postgres:
	docker run --name postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=123321 -p 6432:5432 -d postgres:12-alpine
get_postgres:
	docker exec -it postgres psql -U postgres -d postgres
createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres simple_bank
get_db:
	docker exec -it postgres psql --username=postgres --dbname=postgres simple_bank
dropdb:
	docker exec -it postgres dropdb -U postgres simple_bank

create_migration:
	migrate create -ext sql -dir db/migration -seq init_schema
migrate_up:
	migrate -path db/migration/ -database "postgresql://postgres:123321@localhost:6432/simple_bank?sslmode=disable" --verbose up
migrate_down:
	migrate -path db/migration/ -database "postgresql://postgres:123321@localhost:6432/simple_bank?sslmode=disable" --verbose down

pull_sqlc:
	docker pull sqlc/sqlc
init_sqlc:
	docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc init
generate_sqlc:
	docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc generate

run_test:
	go test -v -cover ./...
clean_test_cache:
	go clean -testcache


.PHONY: postgres,get_postgres,createdb,dropdb
.PHONY: create_migration,migrate_up,migrate_down,migrate_sql
.PHONY: pull_sqlc,init_sqlc,generate_sqlc
.PHONY: run_test,clean_test_cache