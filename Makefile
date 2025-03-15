.PHONY: run build test clean

# Application variables
APP_NAME=todo-api
APP_PORT=8000

# Go variables
GO=go
GOBUILD=$(GO) build
GOTEST=$(GO) test
GOCLEAN=$(GO) clean

# Build and run the application
run:
	@echo "Running $(APP_NAME)..."
	@$(GO) run cmd/api/main.go

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@$(GOBUILD) -o $(APP_NAME) cmd/api/main.go

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@$(GOCLEAN)
	@rm -f $(APP_NAME)

# Run the application with hot reloading using air (if installed)
dev:
	@if command -v air > /dev/null; then \
		echo "Running $(APP_NAME) with hot reloading..."; \
		air; \
	else \
		echo "air is not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		$(MAKE) run; \
	fi

# Show help
help:
	@echo "Available commands:"
	@echo "  make run    - Run the application"
	@echo "  make build  - Build the application"
	@echo "  make test   - Run tests"
	@echo "  make clean  - Clean build artifacts"
	@echo "  make dev    - Run with hot reloading (requires air)"
	@echo "  make help   - Show this help"
