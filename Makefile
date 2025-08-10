include .env
export $(shell sed 's/=.*//' .env)

build:
	go build -C cmd/bot/ -o ../../bin/

run: build
	./bin/bot -env ./.env