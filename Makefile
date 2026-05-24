.PHONY: build run test

build:
	@echo "Building shortix: url shortener service...."
	@go build -o bin/shortix cmd/main.go
	@echo "Build complete: bin/shortix"

run:
	@go run cmd/main.go

test:
	@echo "Running tests with coverage..."
	@go test -v -cover ./...
