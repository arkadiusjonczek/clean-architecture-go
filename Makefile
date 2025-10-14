.PHONY: all
all: lint test

lint:
	golangci-lint run

test:
	go test ./...

run:
	go run ./cmd/server
