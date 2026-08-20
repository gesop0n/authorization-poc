.PHONY: dev start build

dev:
	go run ./cmd/api

build:
	go build -o ./bin/api/ ./cmd/api

start: build
	./bin/api
