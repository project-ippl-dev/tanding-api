# Path to env file
include .env

# Construct database URL
DATABASE_URL='postgres://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable'

# Path to migration files
MIGRATION_PATH=postgresql/migration

# Migration Commands
migrate-up:
	@echo "Running migrations..."
	migrate -path ${MIGRATION_PATH} -database "${DATABASE_URL}" up ${N}

migrate-down:
	@echo "Rolling back migrations..."
	migrate -path ${MIGRATION_PATH} -database "${DATABASE_URL}" down ${N}

migrate-create:
	@echo "Creating new migration file..."
	migrate create -ext sql -dir ${MIGRATION_PATH} ${name}

sqlc-generate:
	sqlc generate

lint:
	golangci-lint run

up:
	docker compose up -d

restart:
	docker compose up -d --force-recreate

gcloud-login:
	gcloud auth application-default login
