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
	go test ./core/... ./common/... ./api/... -v