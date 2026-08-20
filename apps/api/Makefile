include .env
export

MIGRATION_DIR := db/migrations

.PHONY: dev start build

# `make migrate-create create_users` の第2ゴールを migration 名として扱う。
ifeq ($(firstword $(MAKECMDGOALS)),migrate-create)
MIGRATION_NAME := $(word 2,$(MAKECMDGOALS))
ifneq ($(word 2,$(MAKECMDGOALS)),)
.PHONY: $(word 2,$(MAKECMDGOALS))
$(word 2,$(MAKECMDGOALS)):
	@:
endif
endif

migrate-create:
	@test -n "$(MIGRATION_NAME)" || (echo "Usage: make migrate-create <migration-name>"; exit 1)
	goose -dir $(MIGRATION_DIR) create $(MIGRATION_NAME) sql

migrate-up:
	goose -dir $(MIGRATION_DIR) postgres $(DATABASE_URL) up

migrate-down:
	goose -dir $(MIGRATION_DIR) postgres $(DATABASE_URL) down

migrate-status:
	goose -dir $(MIGRATION_DIR) postgres $(DATABASE_URL) status

dev:
	go run ./cmd/api

build:
	go build -o ./bin/api/ ./cmd/api

start: build
	./bin/api
