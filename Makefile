.PHONY: build install install-global test format lint run serve status clean build-linux build-all

# Build the binary with CGO and FTS5 support
build:
	CGO_ENABLED=1 go build -tags "fts5" -o claude_peepee .

# Install to user's bin directory
install: build
	mkdir -p ~/.bin
	cp claude_peepee ~/.bin/

# Install globally
install-global: build
	sudo cp claude_peepee /usr/local/bin/

# Run tests
test:
	CGO_ENABLED=1 go test -tags "fts5" -v ./...

# Format code
format:
	gofmt -w .
	go mod tidy

# Lint code
lint:
	go vet -tags "fts5" ./...
	golangci-lint run --build-tags fts5

# Run the binary
run: build
	./claude_peepee

# Run serve command
serve: build
	./claude_peepee serve

# Run status command
status: build
	./claude_peepee status

# Clean build artifacts
clean:
	rm -f claude_peepee
	rm -f claude_peepee-*

# Build for Linux
build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags "fts5" -o claude_peepee-linux-amd64 .

# Build for all platforms
build-all: build build-linux
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -tags "fts5" -o claude_peepee-darwin-arm64 .
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -tags "fts5" -o claude_peepee-darwin-amd64 .

# Development helpers
dev-remember: build
	./claude_peepee remember $(FACT)

dev-recall: build
	./claude_peepee recall $(QUERY)
