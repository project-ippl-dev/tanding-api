# Path to config file
CONFIG_FILE=config.yaml

# Extract values from YAML
DB_USER=$(shell yq '.database.username' $(CONFIG_FILE))
DB_PASS=$(shell yq '.database.password' $(CONFIG_FILE))
DB_NAME=$(shell yq '.database.name' $(CONFIG_FILE))
DB_HOST=$(shell yq '.database.host' $(CONFIG_FILE))
DB_PORT=$(shell yq '.database.port' $(CONFIG_FILE))

# Construct database URL
DATABASE_URL=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Path to migration files
MIGRATION_PATH=postgresql/migration

# Migration Commands
migrate-up:
	@echo "Running migrations..."
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up ${N}

migrate-down:
	@echo "Rolling back migrations..."
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down ${N}

migrate-create:
	@echo "Creating new migration file..."
	migrate create -ext sql -dir $(MIGRATION_PATH) $(name)

sqlc-generate:
	sqlc generate

lint:
	golangci-lint run

up:
	docker compose up -d

restart:
	docker compose up -d --force-recreate
