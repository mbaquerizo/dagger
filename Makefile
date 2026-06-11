.PHONY: build run test lint migrate seedkey

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

seedkey:
	go run ./cmd/seedkey $(ARGS)
