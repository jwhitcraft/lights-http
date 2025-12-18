# Lights HTTP Server

A simple HTTP server for controlling Govee lights with authentication. 404 and 401 errors redirect to random xkcd comics for entertainment!

## Features

- Turn lights on/off
- Set lights to red, yellow, or orange
- Bearer token authentication
- Configurable via environment variables or .env file

## Endpoints

- `POST /lights/on` - Turn lights on
- `POST /lights/off` - Turn lights off
- `POST /lights/red` - Set lights to red
- `POST /lights/yellow` - Set lights to yellow
- `POST /lights/orange` - Set lights to orange
- `POST /lights/dark-red` - Set lights to dark red
- `POST /lights/rgb` - Set custom RGB color (JSON body: `{"r": 255, "g": 128, "b": 0}`)
- `POST /lights/colortemp` - Set color temperature in Kelvin (JSON body: `{"temperature": 3000}`)
- `POST /lights/brightness` - Set brightness (JSON body: `{"brightness": 50}`)
- `GET /lights/status` - Get status of all devices (returns deviceID, onOff, brightness, color)
- `GET /health` - Health check endpoint (returns overall status, uptime, component checks)
- `GET /ready` - Readiness probe (same as /health)
- `GET /live` - Liveness probe (same as /health)

All endpoints require a Bearer token in the Authorization header.

## Example Usage

Assuming the server is running on `http://localhost:8080` and `BEARER_TOKEN=your-token`:

```bash
# Turn lights on
curl -X POST http://localhost:8080/lights/on \
  -H "Authorization: Bearer your-token"

# Turn lights off
curl -X POST http://localhost:8080/lights/off \
  -H "Authorization: Bearer your-token"

# Set lights to red
curl -X POST http://localhost:8080/lights/red \
  -H "Authorization: Bearer your-token"

# Set lights to yellow
curl -X POST http://localhost:8080/lights/yellow \
  -H "Authorization: Bearer your-token"

# Set lights to orange
curl -X POST http://localhost:8080/lights/orange \
  -H "Authorization: Bearer your-token"

# Set lights to dark red
curl -X POST http://localhost:8080/lights/dark-red \
  -H "Authorization: Bearer your-token"

# Set custom RGB color (orange)
curl -X POST http://localhost:8080/lights/rgb \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{"r": 255, "g": 128, "b": 0}'

# Set color temperature (warm white)
curl -X POST http://localhost:8080/lights/colortemp \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{"temperature": 3000}'

# Set brightness to 50
curl -X POST http://localhost:8080/lights/brightness \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{"brightness": 50}'

# Get device statuses
curl -X GET http://localhost:8080/lights/status \
  -H "Authorization: Bearer your-token"
# Returns: [{"deviceID": "35:CF:DC:6E:00:86:3C:94", "onOff": true, "brightness": 100, "color": {"r": 255, "g": 0, "b": 0}}, ...]

# Get Prometheus metrics (on separate metrics port)
curl -X GET http://localhost:9090/metrics
# Returns Prometheus-formatted metrics for monitoring
```

## Configuration

Set the following environment variables:

- `HOSTNAME` (default: 0.0.0.0)
- `PORT` (default: 8080)
- `METRICS_PORT` (default: 9090)
- `BEARER_TOKEN` (required)
- `GO_ENV` (set to "production" to skip .env loading)

For development, create a `.env` file with the variables.

## Monitoring

The application exposes Prometheus metrics on a separate port for security:

- **API Server**: `http://localhost:8080` (requires authentication)
- **Metrics Server**: `http://localhost:9090/metrics` (no authentication required)

Configure Prometheus to scrape metrics from the metrics endpoint. The metrics include:
- HTTP request counts and latency histograms
- Light operation success/failure counters
- Active connection gauges
- Go runtime metrics

## Development

This project includes a Makefile with common Go development tasks:

```bash
make help          # Show available targets
make build         # Build the binary
make test          # Run tests
make test-verbose  # Run tests with verbose output
make test-cover    # Run tests with coverage report
make run           # Run the application
make fmt           # Format code
make vet           # Run go vet
make lint          # Run golint
make clean         # Clean build artifacts
make deps          # Download dependencies
make check         # Run format, vet, and tests
make docker-build  # Build Docker image
make docker-run    # Run Docker container
```

## Running Natively (Recommended for macOS)

For device discovery to work properly on macOS, run the application natively:

```bash
# Install dependencies
go mod download

# Run the application
go run main.go
```

The application will load configuration from `.env` file or environment variables.

### Troubleshooting Network Issues on macOS

If the container can't discover devices on your local network (e.g., Govee lights), try these options:

1. **Host networking mode** (experimental on macOS):
   ```bash
   docker run --network host \
     -e HOSTNAME=0.0.0.0 \
     -e PORT=8080 \
     -e BEARER_TOKEN=your-secret-token \
     -e GO_ENV=production \
     lights-http
   ```

2. **Privileged mode** (gives full network access):
   ```bash
   docker run --privileged \
     -p 8080:8080 \
     -e HOSTNAME=0.0.0.0 \
     -e PORT=8080 \
     -e BEARER_TOKEN=your-secret-token \
     -e GO_ENV=production \
     lights-http
   ```

**Note**: Docker networking on macOS (Docker Desktop) runs in a VM and may not properly forward UDP broadcasts for device discovery. This issue is specific to macOS; Docker on Linux should work fine. Running natively is recommended for reliable device detection on macOS.

## Running with Docker

## Running

```bash
go run main.go
```

## Testing

```bash
go test ./...
```

## Dependencies

- github.com/swrm-io/go-vee (local copy via replace)
- github.com/joho/godotenv

## Monitoring & Observability

The application includes comprehensive monitoring features:

### Structured Logging
- **JSON format** for easy parsing by log aggregation tools
- **Request IDs** for tracing requests across logs
- **Request/response logging** with timing and metadata
- **Consistent log levels** (DEBUG, INFO, WARN, ERROR)

Example JSON log output:
```json
{
  "time": "2025-12-15T20:20:25.656008-05:00",
  "level": "INFO",
  "source": {
    "function": "main.main",
    "file": "/Users/jwhitcraft/Projects/lights-http/main.go",
    "line": 96
  },
  "msg": "Starting server",
  "addr": "0.0.0.0:8080"
}
```

### Request Tracing
- Unique request IDs generated for each HTTP request
- Request IDs included in response headers (`X-Request-ID`)
- All handler logs include request ID for correlation
- Easy debugging of request flows

### Health Endpoints
- `GET /health` - Overall application health with component status
- `GET /ready` - Kubernetes readiness probe
- `GET /live` - Kubernetes liveness probe

Example health response:
```json
{
  "status": "ok",
  "timestamp": "2025-12-15T20:23:58Z",
  "uptime": "3.6672225s",
  "checks": {
    "controller": {
      "status": "ok",
      "detail": "1 devices connected"
    }
  }
}
```

### Configuration
Set the following environment variables for logging:

- `LOG_LEVEL` (DEBUG, INFO, WARN, ERROR - default: INFO)
- `LOG_FORMAT` (json/text - default: json)

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.