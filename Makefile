.PHONY: build test clean build-all

# Default target
build:
	go build -o bin/analyze-repo cmd/analyze-repo/main.go

# Cross-platform builds
build-all:
	GOOS=windows GOARCH=amd64 go build -o bin/analyze-repo-windows.exe cmd/analyze-repo/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/analyze-repo-darwin cmd/analyze-repo/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/analyze-repo-linux cmd/analyze-repo/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run the tool on current directory
run:
	go run cmd/analyze-repo/main.go .

# Run with verbose output
run-verbose:
	go run cmd/analyze-repo/main.go . --verbose

# Generate JSON output
run-json:
	go run cmd/analyze-repo/main.go . --format=json