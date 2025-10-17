# =========================
# Go commands
# =========================

tidy:
	go mod tidy

run: 
	go run main.go

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