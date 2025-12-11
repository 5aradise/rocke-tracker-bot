include .env
export $(shell sed 's/=.*//' .env)

.PHONY: run build test fmt lint sql-generate migrate-up migrate-down

run: build
	./bin/bot

build:
	go build -C cmd/bot/ -o ../../bin/

test:
	go test ./... -timeout 60s -v -cover -race

fmt:
	go fmt ./...

lint:
	go tool golangci-lint run

sql-generate:
	go tool sqlc generate

migrate-up: 
	GOOSE_MIGRATION_DIR=sql/schema go tool goose sqlite3 ./data.db up

migrate-down: 
	GOOSE_MIGRATION_DIR=sql/schema go tool goose sqlite3 ./data.db down
