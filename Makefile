# =========================
# Go commands
# =========================

tidy:
	go mod tidy

run: 
	go run main.go

# Development mode with hot reload
dev:
	@echo "🔄 Starting development server with hot reload..."
	@echo "📚 Swagger docs will be available at http://localhost:8080/swagger/index.html"
	swag init --parseDependency --parseInternal
	air

build: 
	go build -o main main.go

run-build: build
	./main

test:
	go test ./...

test-integration:
	go test ./tests/integration -v

test-unit:
	go test ./core/... ./support/... ./api/... -v

setup:
	go run main.go setup

seed:
	go run main.go seed

# =========================
# Migration commands (Goose + Atlas)
# =========================

# Run all pending migrations
migrate-up:
	go run main.go migrate

# Rollback last migration
migrate-down:
	go run main.go migrate:down

# Show migration status
migrate-status:
	go run main.go migrate:status

# Reset all migrations (rollback everything)
migrate-reset:
	go run main.go migrate:reset

# ✨ Auto-generate migration from GORM models (Atlas)
# Usage: make migrate-generate name=add_user_age
migrate-generate:
	@if [ -z "$(name)" ]; then \
		echo "❌ Usage: make migrate-generate name=<migration_name>"; \
		exit 1; \
	fi
	go run main.go migrate:generate $(name)

# Create empty migration file for manual SQL
# Usage: make migrate-create name=custom_change
migrate-create:
	@if [ -z "$(name)" ]; then \
		read -p "Enter migration name: " name; \
		go run main.go migrate:create $$name sql; \
	else \
		go run main.go migrate:create $(name) sql; \
	fi

# Re-hash migrations (run after manual edits to migration files)
migrate-hash:
	atlas migrate hash --env local

# =========================
# Swagger commands
# =========================

swagger:
	swag init --parseDependency --parseInternal

swagger-fmt:
	swag fmt

# =========================
# Docker commands
# =========================

up:
	docker-compose up -d
	@echo "Containers started in detached mode. Use 'make logs' to follow logs."

down:
	docker-compose down

logs:
	docker-compose logs -f

rebuild:
	docker-compose down -v
	docker-compose up --build -d
	@echo "Containers rebuilt and started. Use 'make logs' to follow logs."