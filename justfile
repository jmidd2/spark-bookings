# justfile for Go project with Air hot reloading

# Default recipe - shows available commands
default:
    @just --list

# Install Air if not already installed
install-air:
    #!/usr/bin/env bash
    if ! command -v air &> /dev/null; then
        echo "Installing Air..."
        go install github.com/air-verse/air@latest
        echo "Air installed successfully!"
    else
        echo "Air is already installed"
    fi

install-lint:
    #!/usr/bin/env bash
    if ! command -v golangci-lint &> /dev/null; then
        echo "Installing golangci-lint..."
        go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
        echo "golangci-lint installed successfully!"
    else
        echo "golangci-lint is already installed"
    fi

# Start development server with hot reloading
dev: install-air
    @echo "Starting development server with hot reloading..."
    @echo "Your app will run on :8080, but visit http://localhost:8090 for live reload!"
    APP_ENV=development air

# Start development with verbose logging
dev-verbose: install-air
    @echo "Starting development server with verbose logging..."
    @echo "Your app will run on :8080, but visit http://localhost:8090 for live reload!"
    APP_ENV=development air -d

# Initialize Air configuration file if it doesn't exist
init-air: install-air
    #!/usr/bin/env bash
    if [ ! -f .air.toml ]; then
        echo "Creating Air configuration file..."
        air init
        echo "Air configuration created at .air.toml"
    else
        echo "Air configuration already exists"
    fi

# Start development with custom Air config
dev-config CONFIG: install-air
    @echo "Starting development server with config: {{CONFIG}}"
    air -c {{CONFIG}}

# Clean build artifacts
clean:
    @echo "Cleaning build artifacts..."
    rm -rf tmp/
    go clean

# Build the project
build:
    @echo "Building project..."
    go build -o bin/booking-display .

# Run tests
test:
    @echo "Running tests..."
    go test ./...

# Run tests with coverage
test-coverage:
    @echo "Running tests with coverage..."
    go test -cover ./...

# Format code
fmt:
    @echo "Formatting code..."
    go fmt ./...

# Run linter (requires golangci-lint)
lint: install-lint
    @echo "Running linter..."
    golangci-lint run

# Full development setup - initializes Air config and starts dev server
setup: init-air install-lint dev

# Show Air version
air-version: install-air
    air -v

# Open browser to the live reload URL (macOS)
open:
    @echo "Opening browser to live reload URL..."
    open http://localhost:8090

# Open browser to the live reload URL (Linux)
open-linux:
    @echo "Opening browser to live reload URL..."
    xdg-open http://localhost:8090