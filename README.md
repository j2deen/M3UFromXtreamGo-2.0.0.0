# M3UFromXtream

This was completely adopted from https://github.com/BillOatmanWork/M3UFromXtream/tree/V1.0.0.0 thank you for your original code I added some capabilities for my purpose that might be useful to others so I've extended it and converted it into GoLang so it's a bit more selfcontained for docker. I have kept the same license and suggested charitable donation as a thank you to the orginal author.

A dual-mode application to create M3U playlists from Xtream Code API. Supports both CLI and web server modes, with Docker containerization for easy deployment.

## Features

- **CLI Mode**: Traditional command-line interface for generating M3U files
- **Web Server Mode**: HTTP server with REST endpoints for on-demand playlist generation
- **Flexible Configuration**: Support for environment variables, config files, or CLI arguments
- **Docker Support**: Containerized deployment with minimal image size (17.6MB)
- **Static Binary**: Self-contained executable with no external dependencies (CGO_ENABLED=0)
- **Health Checks**: Built-in health endpoint for container orchestration
- **Comprehensive Logging**: Configurable log levels (DEBUG, INFO, WARN, ERROR) with detailed output
- **Error Handling**: Partial failure recovery, detailed error messages, and authentication feedback

## Quick Start

### CLI Mode (Traditional Usage)

```bash
./m3ufromxtream http://example.com:8080 user pass output.m3u
```

### Web Server Mode

Using environment variables:
```bash
export M3U_XTREAM_BASE_URL=http://example.com:8080
export M3U_XTREAM_USERNAME=your_username
export M3U_XTREAM_PASSWORD=your_password
./m3ufromxtream
```

Then access the M3U playlist at: `http://localhost:8080/m3u`

### Docker

Build and run:
```bash
docker build -t m3ufromxtream:latest .

docker run -d -p 8080:8080 \
  -e M3U_XTREAM_BASE_URL=http://example.com:8080 \
  -e M3U_XTREAM_USERNAME=your_username \
  -e M3U_XTREAM_PASSWORD=your_password \
  m3ufromxtream:latest
```

Or use the Makefile:
```bash
make docker
make docker-run
```

## Configuration

### Environment Variables

- `M3U_MODE` - Operating mode: "web" (default) or "cli"
- `M3U_SERVER_PORT` - Server port (default: 8080)
- `M3U_SERVER_HOST` - Bind address (default: 0.0.0.0)
- `M3U_XTREAM_BASE_URL` - Xtream API URL (required)
- `M3U_XTREAM_USERNAME` - Xtream username (required)
- `M3U_XTREAM_PASSWORD` - Xtream password (required)
- `M3U_XTREAM_REQUEST_TIMEOUT` - API request timeout in seconds (default: 30)
- `M3U_CONFIG_FILE` - Path to config file (default: ./config.json)
- `M3U_LOG_LEVEL` - Logging level: DEBUG, INFO, WARN, ERROR (default: INFO)

### Configuration File

Create a `config.json` file:
```json
{
  "mode": "web",
  "server": {
    "port": 8080,
    "host": "0.0.0.0"
  },
  "xtream": {
    "base_url": "http://example.com:8080",
    "username": "your_username",
    "password": "your_password"
  }
}
```

See `config.example.json` for a complete configuration template.

### Configuration Priority

Configuration is loaded with the following priority (highest to lowest):
1. CLI arguments (CLI mode only)
2. Environment variables
3. Configuration file
4. Default values

## Web Server Endpoints

- **GET /m3u** - Returns M3U playlist
  - Content-Type: `application/vnd.apple.mpegurl`
  - Downloads as `playlist.m3u`

- **GET /health** - Health check endpoint
  - Returns: `{"status":"healthy","version":"1.0.0.0","timestamp":"..."}`

- **GET /config** - Configuration info (debug endpoint)
  - Returns current configuration with password redacted

## Building

### Using Makefile

```bash
make build          # Build for current platform
make build-linux    # Build static binary for Linux
make build-all      # Build for all platforms
make docker         # Build Docker image
make clean          # Remove build artifacts
```

### Manual Build

Build static binary (Linux):
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o m3ufromxtream .
```

Build for current platform:
```bash
CGO_ENABLED=0 go build -o m3ufromxtream .
```

## Logging and Error Handling

### Logging Levels

The application supports four logging levels:

- **DEBUG** - Detailed diagnostic information (API calls, response sizes, category processing)
- **INFO** - General informational messages (startup, configuration, progress)
- **WARN** - Warning messages (partial failures, skipped items)
- **ERROR** - Error messages (API failures, connection issues)

Set the log level via environment variable:
```bash
M3U_LOG_LEVEL=DEBUG ./m3ufromxtream
```

Or in the configuration file:
```json
{
  "logging": {
    "level": "DEBUG"
  }
}
```

### Error Handling Features

- **Partial Failure Recovery** - If some categories fail to load, the application continues processing other categories and logs warnings
- **Detailed Error Messages** - Authentication failures, connection timeouts, and JSON parsing errors include helpful context
- **Password Sanitization** - Passwords are redacted from log output for security
- **Request Logging** - HTTP requests include method, path, status code, and duration
- **Progress Tracking** - Category and stream processing progress is logged

### Example Log Output

```
[2026-01-07 20:37:42] INFO - M3UFromXtream starting - Version 1.0.0.0
[2026-01-07 20:37:42] INFO - Running in CLI mode
[2026-01-07 20:37:42] INFO - Starting M3U generation from Xtream API
[2026-01-07 20:37:42] INFO - Fetching categories from Xtream API...
[2026-01-07 20:37:42] INFO - Successfully fetched 25 categories
[2026-01-07 20:37:42] DEBUG - Processing category 1/25: Sports (ID: 123)
[2026-01-07 20:37:42] INFO - Category 'Sports': processed 45 streams
[2026-01-07 20:37:42] WARN - Failed to fetch streams for category 'Premium' (ID: 456): HTTP 401
[2026-01-07 20:37:42] INFO - M3U generation complete - Total streams: 450, Skipped: 5, Failed categories: 1
```

## Use Cases

### Emby Scheduled Task

1. Deploy the Docker container
2. Configure Emby to fetch the M3U playlist from: `http://container-ip:8080/m3u`
3. Set up a scheduled task in Emby to refresh the playlist periodically

### Media Server Integration

The web interface allows any media server to fetch updated playlists on-demand without manual regeneration.

### Debugging Issues

Enable DEBUG logging to troubleshoot connection or API issues:
```bash
docker run -p 8080:8080 \
  -e M3U_LOG_LEVEL=DEBUG \
  -e M3U_XTREAM_BASE_URL=http://api.example.com:8080 \
  -e M3U_XTREAM_USERNAME=user \
  -e M3U_XTREAM_PASSWORD=pass \
  m3ufromxtream:latest
```

## Technical Details

- **Language**: Go 1.21+
- **Dependencies**: Standard library only (net/http, encoding/json)
- **Binary Size**: ~7.7MB (unstripped), ~5.5MB (stripped)
- **Docker Image**: 17.6MB (Alpine-based)
- **Static Linking**: CGO_ENABLED=0 ensures portability across Linux distributions

## Good Karma

This application is free. If you find it of value and have the means, please consider making a donation to a local charity that benefits children.



