.PHONY: all
all: run

test:
	go test ./...

run:
	go run ./cmd/server
