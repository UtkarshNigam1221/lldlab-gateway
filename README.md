# lldlab-gateway

API Gateway for all frontend requests - a high-performance reverse proxy built with Go that routes requests to backend microservices.

## Features

- **Request Proxying**: Forward requests to multiple backend services
- **CORS Support**: Configurable CORS policies for frontend integration
- **Middleware Chain**: Logging, recovery, and authentication middleware
- **Health Checks**: Built-in health check endpoint
- **Graceful Shutdown**: Proper cleanup on termination signals
- **Configuration-Driven**: YAML-based route and service configuration
- **Docker Support**: Production-ready containerization

## Project Structure

```
lldlab-gateway/
├── cmd/
│   └── gateway/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration loader
│   ├── handler/
│   │   └── health.go         # Health check handler
│   ├── middleware/
│   │   ├── auth.go           # Authentication middleware
│   │   ├── logging.go        # Request logging
│   │   └── recovery.go       # Panic recovery
│   └── proxy/
│       └── proxy.go          # Reverse proxy handler
├── config.yaml               # Service configuration
├── Dockerfile                # Container image
├── docker-compose.yaml       # Docker Compose setup
├── Makefile                  # Build automation
└── go.mod                    # Go dependencies
```

## Quick Start

### Prerequisites

- Go 1.21 or later
- Docker (optional)

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd lldlab-gateway
```

2. **Install dependencies**
```bash
make deps
```

3. **Configure routes**

Edit `config.yaml` to define your backend services:
```yaml
routes:
  - path: "/api/users"
    prefix: true
    target: "http://localhost:8081"
    methods: ["GET", "POST", "PUT", "DELETE"]
```

4. **Run the gateway**
```bash
make run
```

The gateway will start on `http://localhost:8080`

## Configuration

### Server Settings

```yaml
server:
  port: 8080
  host: "0.0.0.0"
  read_timeout: 15s
  write_timeout: 15s
  idle_timeout: 60s
```

### CORS Configuration

```yaml
cors:
  allowed_origins: ["*"]
  allowed_methods: ["GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"]
  allowed_headers: ["Content-Type", "Authorization"]
  allow_credentials: true
  max_age: 3600
```

### Route Definitions

- **path**: URL path to match
- **prefix**: If true, matches all paths starting with this prefix
- **target**: Backend service URL
- **methods**: Allowed HTTP methods

## Development

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Format Code

```bash
make fmt
```

## Docker Deployment

### Build Image

```bash
make docker-build
```

### Run with Docker Compose

```bash
make docker-run
```

### Stop Services

```bash
make docker-stop
```

## Endpoints

### Health Check

```bash
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2026-01-20T19:19:00Z",
  "service": "api-gateway"
}
```

### Proxied Routes

All configured routes are automatically proxied to their respective backend services.

## Middleware

1. **Recovery**: Catches panics and returns 500 errors
2. **Logging**: Logs all requests with timing and status information
3. **Auth**: Optional authentication layer (currently pass-through)
4. **CORS**: Handles cross-origin requests

## Environment Variables

- `CONFIG_PATH`: Path to configuration file (default: `config.yaml`)

## License

MIT
