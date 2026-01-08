# Stage 1: Build with static linking
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git for any potential dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build with static linking (CGO_ENABLED=0)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o m3ufromxtream .

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests and wget for health checks
RUN apk add --no-cache ca-certificates wget

# Copy binary from builder stage
COPY --from=builder /app/m3ufromxtream .

# Copy example config (can be overridden with volume mount)
COPY config.example.json ./config.json

# Expose default port
EXPOSE 8080

# Health check endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default to web server mode
ENV M3U_MODE=web

# Run the application
ENTRYPOINT ["./m3ufromxtream"]
