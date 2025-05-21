# Path to env file
#include .env ## WARNING: make sure your env is correct

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

mock-generate:
	mockgen \
          -destination=./mocks/db/db_mock.go \
          -package=mock_db \
          -source=./internal/db/db.go DBTX
	go generate ./...; \
    go mod tidy

unit-test:
	@PACKAGES="$$(go list ./... | grep -Ev '/(mocks|testutils|cmd|config|postgresql|db)(/|$$)')"; \
	go test $$PACKAGES \
	  -race \
	  -coverpkg=$$COVERPKGS \
      -covermode=atomic \
      -coverprofile=coverage.out

cover:
	go tool cover -func coverage.out