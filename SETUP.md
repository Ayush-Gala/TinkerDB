# TinkerDB Setup Guide

Quick setup guide for getting TinkerDB up and running. Please note that these steps set up the database to work as I intended it to. If you wish to tinker with the code (pun intended) then feel free to set the project up as per your requirements.

## Prerequisites

### 1. Install Go

**Verify:**
```bash
go version  # Should show Go 1.21 or later
```

### 2. Install Protocol Buffers Compiler

**Verify:**
```bash
protoc --version  # Should show libprotoc 3.x or later
```

## Setup Steps

### 1. Navigate to Project Directory
```bash
cd TinkerDB
```

### 2. Install Go Tools and Dependencies
```bash
make deps
```

This installs:
- `protoc-gen-go` - Protocol Buffers Go plugin
- `protoc-gen-go-grpc` - gRPC Go plugin
- All Go dependencies from `go.mod`

### 3. Generate Protocol Buffers Code
```bash
make proto
```

This generates:
- `proto/kvstore/kvstore.pb.go` - Protocol Buffers messages
- `proto/kvstore/kvstore_grpc.pb.go` - gRPC service code

### 4. Run Tests (Optional but Recommended)
```bash
make test
```
Expected output: All tests pass with coverage report.

### 5. Build the Server
```bash
make build
```

This creates: `bin/tinkerdb-server`

## Running TinkerDB

### Start the Server

**Option 1: Using Make**
```bash
make server
```

**Option 2: Using Go**
```bash
go run cmd/server/main.go
```

**To use a custom Port:**
```bash
TINKERDB_PORT=<port number> make server
```

### Expected Output
```
TinkerDB server starting on port 50051...
Server is ready to accept connections
```

## Troubleshooting

### Problem: `protoc: command not found`
**Solution:** Install Protocol Buffers compiler (see Prerequisites)

### Problem: `cannot find package "google.golang.org/grpc"`
**Solution:**
```bash
go mod download
go mod tidy
```

### Problem: `*.pb.go files not found`
**Solution:**
```bash
make proto
```

### Problem: Server port already in use
**Solution:**
```bash
# Find and kill process using port 50051
lsof -ti:50051 | xargs kill -9

# Or use a different port
TINKERDB_PORT=8080 make server
```

## System Requirements

- **Go**: 1.21 or later
- **Protocol Buffers**: 3.x or later
- **OS**: macOS, Linux, or Windows
- **Memory**: 100MB minimum
- **Disk**: 50MB for dependencies and build artifacts

## Testing

### Running Tests

**All tests:**
```bash
make test
```

**With coverage report:**
```bash
make test-coverage
open coverage.html
```

**Specific package:**
```bash
go test -v ./internal/storage/
go test -v ./internal/server/
go test -v ./test/
```

**Run benchmarks:**
```bash
go test -bench=. -benchmem ./...
```
