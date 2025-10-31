.PHONY: build test clean install

# Build the binary
build:
	go build -o leakyrepo .

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Install to /usr/local/bin (requires sudo)
install: build
	sudo cp leakyrepo /usr/local/bin/

# Clean build artifacts
clean:
	rm -f leakyrepo
	rm -f coverage.out

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

