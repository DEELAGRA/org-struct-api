include .env
export
export DB_USER DB_PASSWORD DB_HOST DB_PORT DB_NAME DB_SSLMODE
export PROJECT_ROOT=$(shell pwd)
DATABASE_URL ?= "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

env-up:
	@docker compose up -d org-struct-api-postgres
env-down:
	@docker compose down org-struct-api-postgres
env-cleanup:
	@ read -p "Clear all volume environment files? Risk of data loss. [Y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres && \
		sudo rm -rf out/pgdata && \
		echo "Environment files deleted"; \
	else \
		echo "File deletion canceled"; \
	fi
run:
	@go run cmd/main.go
 

migrate-up:
	@make migrate-action action=up
migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "The required parameter 'action' is missing. Example: make migrate-action action=up/down"; \
		exit 1; \
	fi; \
	goose -dir db/migrations postgres $(DATABASE_URL) "$(action)"

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "The required parameter 'seq' is missing. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	mkdir -p db/migrations; \
	goose -dir db/migrations create $$seq sql