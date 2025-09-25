.PHONY: build run dev clean test lint docker-up docker-down docker-logs docker-clean

DOCKERFILE_PATH = deployments/

# Build the application
build:
	go build -o bin/shareVault ./cmd/shareVault

# Run the application
run:
	go run ./cmd/shareVault

# Development mode with hot reload (requires air)
dev: 
	~/go/bin/air -c .air.toml

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Lint code
lint:
	golangci-lint run --no-config

# Docker commands
docker-up:
	docker compose -f $(DOCKERFILE_PATH)docker-compose.yml up -d

docker-down:
	docker compose -f $(DOCKERFILE_PATH)docker-compose.yml down

docker-logs:
	docker compose -f $(DOCKERFILE_PATH)docker-compose.yml logs -f

docker-clean:
	docker compose -f $(DOCKERFILE_PATH)docker-compose.yml down -v --remove-orphans