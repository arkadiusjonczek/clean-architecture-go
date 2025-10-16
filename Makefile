.PHONY: all
all: lint test

lint:
	golangci-lint run

test:
	go test ./...

generate:
	go generate ./...

run:
	go run ./cmd/server
