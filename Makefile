include .env
export $(shell sed 's/=.*//' .env)

.PHONY: build run migrate-up migrate-down

run: build
	./bin/bot

build:
	go build -C cmd/bot/ -o ../../bin/

migrate-up: 
	GOOSE_MIGRATION_DIR=sql/schema goose sqlite3 ./data.db up

migrate-down: 
	GOOSE_MIGRATION_DIR=sql/schema goose sqlite3 ./data.db down
	