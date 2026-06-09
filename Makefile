.PHONY: build run test lint

build:
	go build ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./...

lint:
	go vet ./...