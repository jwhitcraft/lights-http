# Makefile for lights-http project
# Follows Go project conventions

.PHONY: help build test test-verbose test-cover clean run fmt vet lint deps docker-build docker-run

# Variables
BINARY_NAME=lights-http
BINARY_PATH=bin/$(BINARY_NAME)
DOCKER_IMAGE=lights-http
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*')

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build the binary
build: ## Build the Go binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_PATH) .

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	@go test ./...

# Run tests with verbose output
test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@go test -v ./...

# Run tests with coverage
test-cover: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Run the application
run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	@go run .

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Run golint (if available)
lint: ## Run golint
	@echo "Running golint..."
	@golint ./... || echo "golint not installed, run: go install golang.org/x/lint/golint@latest"

# Download dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Build Docker image
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

# Build Docker image for linux/amd64
docker-build-amd64: ## Build Docker image for linux/amd64
	@echo "Building Docker image for linux/amd64..."
	@docker build --platform linux/amd64 -t $(DOCKER_IMAGE) .

# Run Docker container
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 8080:8080 \
		-e HOSTNAME=0.0.0.0 \
		-e PORT=8080 \
		-e BEARER_TOKEN=your-token \
		-e GO_ENV=development \
		$(DOCKER_IMAGE)

# Run all checks (format, vet, test)
check: fmt vet test ## Run format, vet, and tests

# Install development tools
install-tools: ## Install development tools (golint, etc.)
	@echo "Installing development tools..."
	@go install golang.org/x/lint/golint@latest