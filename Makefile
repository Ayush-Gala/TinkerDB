.PHONY: proto server client test clean

# Generate protobuf and gRPC code
proto:
	@echo "Generating protobuf code..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/kvstore.proto

# Run the server
server:
	@echo "Starting TinkerDB server..."
	@go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the server binary
build:
	@echo "Building TinkerDB server..."
	@go build -o bin/tinkerdb-server cmd/server/main.go
	@echo "Server binary created: bin/tinkerdb-server"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Run all checks (format, vet, test)
check: test
	@echo "Running go fmt..."
	@go fmt ./...
	@echo "Running go vet..."
	@go vet ./...
	@echo "All checks passed!"

