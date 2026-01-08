# Docker Deployment Guide - M3UFromXtream

This guide explains how to deploy and run the M3UFromXtream Docker container on another machine.

## Option 1: Using Pre-Built Docker Image (Fastest)

If you have the `m3ufromxtream.tar` file:

### 1. Transfer the Image

Copy the tar file to your target machine using one of these methods:

```bash
# Using scp
scp m3ufromxtream.tar user@target-machine:/path/to/destination/

# Using rsync
rsync -avz m3ufromxtream.tar user@target-machine:/path/to/destination/

# Or use USB drive, network share, etc.
```

### 2. Load the Image on Target Machine

```bash
docker load -i m3ufromxtream.tar
```

Verify the image loaded:
```bash
docker images | grep m3ufromxtream
```

### 3. Run the Container

```bash
docker run -d \
  --name m3ufromxtream \
  -p 8080:8080 \
  -e M3U_XTREAM_BASE_URL=http://your-xtream-server:8080 \
  -e M3U_XTREAM_USERNAME=your_username \
  -e M3U_XTREAM_PASSWORD=your_password \
  --restart unless-stopped \
  m3ufromxtream:latest
```

## Option 2: Building from Source on Target Machine

### 1. Transfer the Source Code

```bash
# Create a tarball of the source
tar -czf m3ufromxtream-source.tar.gz \
  Dockerfile \
  *.go \
  go.mod \
  config.example.json \
  Makefile

# Transfer to target machine
scp m3ufromxtream-source.tar.gz user@target-machine:/path/to/destination/
```

### 2. Extract and Build on Target Machine

```bash
# Extract
tar -xzf m3ufromxtream-source.tar.gz
cd m3ufromxtream

# Build the image
docker build -t m3ufromxtream:latest .
```

### 3. Run the Container

Same as Option 1, Step 3 above.

## Configuration

### Required Environment Variables

You must provide these when running the container:

- `M3U_XTREAM_BASE_URL` - Your Xtream API server URL (e.g., http://provider.com:8080)
- `M3U_XTREAM_USERNAME` - Your Xtream username
- `M3U_XTREAM_PASSWORD` - Your Xtream password

### Optional Environment Variables

- `M3U_SERVER_PORT` - Port to run on (default: 8080)
- `M3U_LOG_LEVEL` - Logging level: DEBUG, INFO, WARN, ERROR (default: INFO)
- `M3U_XTREAM_REQUEST_TIMEOUT` - API timeout in seconds (default: 30)

### Using a Config File (Alternative to Environment Variables)

Create a `config.json` file:

```json
{
  "mode": "web",
  "server": {
    "port": 8080,
    "host": "0.0.0.0"
  },
  "xtream": {
    "base_url": "http://your-server:8080",
    "username": "your_username",
    "password": "your_password"
  }
}
```

Mount it when running:

```bash
docker run -d \
  --name m3ufromxtream \
  -p 8080:8080 \
  -v /path/to/config.json:/app/config.json \
  --restart unless-stopped \
  m3ufromxtream:latest
```

## Accessing the Service

Once running, access the M3U playlist at:

```
http://your-machine-ip:8080/m3u
```

Other endpoints:
- Health check: `http://your-machine-ip:8080/health`
- Config info: `http://your-machine-ip:8080/config`

## Container Management

### View Logs

```bash
# View recent logs
docker logs m3ufromxtream

# Follow logs in real-time
docker logs -f m3ufromxtream

# View last 50 lines
docker logs --tail 50 m3ufromxtream
```

### Stop Container

```bash
docker stop m3ufromxtream
```

### Start Container

```bash
docker start m3ufromxtream
```

### Restart Container

```bash
docker restart m3ufromxtream
```

### Remove Container

```bash
docker stop m3ufromxtream
docker rm m3ufromxtream
```

### Update Container

```bash
# Stop and remove old container
docker stop m3ufromxtream
docker rm m3ufromxtream

# Load new image or rebuild
docker load -i m3ufromxtream-new.tar
# or
docker build -t m3ufromxtream:latest .

# Run with same configuration
docker run -d --name m3ufromxtream -p 8080:8080 ...
```

## Troubleshooting

### Container Won't Start

Check logs:
```bash
docker logs m3ufromxtream
```

Verify image exists:
```bash
docker images | grep m3ufromxtream
```

### Can't Access the Service

Check if container is running:
```bash
docker ps | grep m3ufromxtream
```

Check port binding:
```bash
docker port m3ufromxtream
```

Test from within the container:
```bash
docker exec m3ufromxtream wget -O- http://localhost:8080/health
```

### Enable Debug Logging

Recreate container with DEBUG logging:
```bash
docker stop m3ufromxtream
docker rm m3ufromxtream

docker run -d \
  --name m3ufromxtream \
  -p 8080:8080 \
  -e M3U_LOG_LEVEL=DEBUG \
  -e M3U_XTREAM_BASE_URL=http://your-server:8080 \
  -e M3U_XTREAM_USERNAME=your_username \
  -e M3U_XTREAM_PASSWORD=your_password \
  m3ufromxtream:latest
```

Then check logs for detailed information.

### Authentication Errors

Verify your credentials by checking the config endpoint:
```bash
curl http://localhost:8080/config
```

The password will be redacted but you can verify other settings.

## Using with Emby/Jellyfin/Plex

1. Ensure the container is running and accessible
2. In your media server, add an M3U playlist source
3. Use the URL: `http://docker-host-ip:8080/m3u`
4. Set up automatic refresh (recommended: every 24 hours)

## Docker Compose (Optional)

Create a `docker-compose.yml` for easier management:

```yaml
version: '3.8'

services:
  m3ufromxtream:
    image: m3ufromxtream:latest
    container_name: m3ufromxtream
    ports:
      - "8080:8080"
    environment:
      - M3U_XTREAM_BASE_URL=http://your-server:8080
      - M3U_XTREAM_USERNAME=your_username
      - M3U_XTREAM_PASSWORD=your_password
      - M3U_LOG_LEVEL=INFO
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
```

Run with:
```bash
docker-compose up -d
```

## Security Notes

- Keep your credentials secure - use environment variables or mounted config files
- If exposing to the internet, consider placing behind a reverse proxy with authentication
- Regularly update the container for security patches

## System Requirements

- Docker 20.10+ or compatible runtime
- 20MB disk space for the image
- Minimal CPU/RAM (idle: <10MB RAM)
- Network access to your Xtream API provider
