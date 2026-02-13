# Install development tools
tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0
	go install github.com/evilmartians/lefthook/v2@v2.1.0

# Setup git hooks
hooks:
	lefthook install

# Run tests
test:
	go test -v -race -count=1 ./...

# Run linter
lint:
	golangci-lint run ./...

# Run linter and tests
check: lint test

# Setup development environment
setup: tools hooks
