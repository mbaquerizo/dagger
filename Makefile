.PHONY: build run test lint migrate

build:
	go build ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./...

lint:
	go vet ./...

migrate:
	go run ./cmd/migrate