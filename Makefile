include .env
export

export PROJECT_ROOT=$(shell pwd)


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