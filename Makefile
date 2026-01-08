.PHONY: build build-linux build-all docker docker-run test clean help

VERSION := 1.0.0.0
BINARY_NAME := m3ufromxtream

# Build for current platform
build:
	@echo "Building for current platform..."
	@CGO_ENABLED=0 go build -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

# Build for Linux with optimizations (static binary)
build-linux:
	@echo "Building for Linux (static binary)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -ldflags="-s -w" -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "Building for all platforms..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -ldflags="-s -w" -o $(BINARY_NAME)-linux-amd64 .
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
		go build -ldflags="-s -w" -o $(BINARY_NAME)-linux-arm64 .
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
		go build -ldflags="-s -w" -o $(BINARY_NAME)-darwin-amd64 .
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
		go build -ldflags="-s -w" -o $(BINARY_NAME)-darwin-arm64 .
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
		go build -ldflags="-s -w" -o $(BINARY_NAME)-windows-amd64.exe .
	@echo "Build complete for all platforms"

# Build Docker image
docker:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):$(VERSION) .
	@docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest
	@echo "Docker image built: $(BINARY_NAME):$(VERSION)"

# Run Docker container with example environment variables
docker-run:
	@echo "Running Docker container..."
	@docker run -d -p 8080:8080 \
		--name $(BINARY_NAME) \
		-e M3U_XTREAM_BASE_URL=http://example.com:8080 \
		-e M3U_XTREAM_USERNAME=user \
		-e M3U_XTREAM_PASSWORD=pass \
		$(BINARY_NAME):latest
	@echo "Container started. Access endpoints at:"
	@echo "  http://localhost:8080/m3u"
	@echo "  http://localhost:8080/health"
	@echo ""
	@echo "To view logs: docker logs -f $(BINARY_NAME)"
	@echo "To stop: docker stop $(BINARY_NAME)"
	@echo "To remove: docker rm $(BINARY_NAME)"

# Stop and remove Docker container
docker-stop:
	@echo "Stopping and removing Docker container..."
	@docker stop $(BINARY_NAME) 2>/dev/null || true
	@docker rm $(BINARY_NAME) 2>/dev/null || true
	@echo "Container stopped and removed"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME)-*
	@rm -f *.exe
	@rm -f *.m3u
	@echo "Clean complete"

# Display help
help:
	@echo "M3UFromXtream Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build         - Build for current platform"
	@echo "  build-linux   - Build static binary for Linux (CGO_ENABLED=0)"
	@echo "  build-all     - Build for all platforms (Linux, macOS, Windows)"
	@echo "  docker        - Build Docker image"
	@echo "  docker-run    - Run Docker container with example config"
	@echo "  docker-stop   - Stop and remove Docker container"
	@echo "  test          - Run tests"
	@echo "  clean         - Remove build artifacts"
	@echo "  help          - Display this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build-linux              # Build static binary for Linux"
	@echo "  make docker                    # Build Docker image"
	@echo "  make docker-run                # Run in Docker"
	@echo ""
	@echo "CLI Usage:"
	@echo "  ./$(BINARY_NAME) <url> <user> <pass> [output.m3u]"
	@echo ""
	@echo "Web Server Usage:"
	@echo "  M3U_XTREAM_BASE_URL=http://... M3U_XTREAM_USERNAME=user \\"
	@echo "  M3U_XTREAM_PASSWORD=pass ./$(BINARY_NAME)"
