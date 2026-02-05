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

migrate:
	go run main.go migrate

seed:
	go run main.go seed

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