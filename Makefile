.PHONY: help api-start api-test frontend-start frontend-test frontend-lint frontend-build all-tests clean

# Default target
help:
	@echo "SplitFlap Project Commands"
	@echo ""
	@echo "API Commands:"
	@echo "  make api-start       - Start the Spring Boot API server"
	@echo "  make api-test        - Run API tests"
	@echo ""
	@echo "Frontend Commands:"
	@echo "  make frontend-start  - Start the SolidJS development server"
	@echo "  make frontend-test   - Run frontend tests"
	@echo "  make frontend-build  - Build frontend for production"
	@echo ""
	@echo "Combined Commands:"
	@echo "  make all-tests       - Run all tests (API + frontend)"
	@echo "  make clean           - Clean build artifacts"

# API Commands
api-start:
	@echo "Starting API server..."
	cd splitflap-api-go && go run ./cmd/api

api-test:
	@echo "Running API tests..."
	cd splitflap-api-go && go test ./...

# Frontend Commands
frontend-start:
	@echo "Starting frontend development server..."
	cd splitflap-web-solid && pnpm dev

frontend-test:
	@echo "Running frontend tests..."
	cd splitflap-web-solid && pnpm test

frontend-build:
	@echo "Building frontend for production..."
	cd splitflap-web-solid && pnpm build

# Combined Commands
all-tests: api-test frontend-test
	@echo "All tests completed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	cd splitflap-api-go && go clean -testcache
	cd splitflap-web-solid && rm -rf dist node_modules/.vite
	@echo "Clean complete!"
