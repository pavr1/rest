# 🍺 Bar Restaurant Management System

A microservices-based restaurant and bar management application built with Go and Docker.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        UI SERVICE                           │
│                     (Port 3000)                             │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    GATEWAY SERVICE                          │
│                     (Port 8082)                             │
│         Request Routing • Auth • Health Monitoring          │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                   BUSINESS SERVICES                         │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐  │
│  │ Session Service │  │  Menu Service   │  │   Orders    │  │
│  │   (Port 8087)   │  │   (Planned)     │  │  (Planned)  │  │
│  └────────┬────────┘  └────────┬────────┘  └──────┬──────┘  │
└───────────┼────────────────────┼─────────────────┼──────────┘
            │                    │                 │
┌───────────▼────────────────────▼─────────────────▼──────────┐
│                     DATA SERVICE                            │
│                     (Port 8086)                             │
│            Configuration • Database Access                  │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                      POSTGRESQL                             │
│                     (Port 5432)                             │
└─────────────────────────────────────────────────────────────┘
```

## Quick Start

```bash
# Start all services
make fresh

# Open UI
open http://localhost:3000
```

## Services

| Service | Port | Description | Documentation |
|---------|------|-------------|---------------|
| **UI Service** | 3000 | Frontend web application | [ui-service/](ui-service/) |
| **Gateway Service** | 8082 | API Gateway, routing, auth | [gateway-service/](gateway-service/) |
| **Session Service** | 8087 | Authentication & JWT | [session-service/README.md](session-service/README.md) |
| **Data Service** | 8086 | Config & database access | [data-service/README.md](data-service/README.md) |
| **PostgreSQL** | 5432 | Database | Managed by data-service |

## Commands

### Service Management

```bash
make start    # Start all services
make stop     # Stop all services
make restart  # Restart all services
make status   # Show status of all services
make fresh    # Clean install everything
make clean    # Remove all containers
```

### Testing

```bash
make test           # Run all tests
make test-data      # Test data-service
make test-session   # Test session-service
make test-gateway   # Test gateway-service
```

### Logs

```bash
make logs s=data      # View data-service logs
make logs s=session   # View session-service logs
make logs s=gateway   # View gateway-service logs
make logs s=ui        # View ui-service logs
```

## Project Structure

```
bar-restaurant/
├── Makefile              # Global orchestration
├── README.md             # This file
├── docs/                 # Documentation
│   └── 1. architecture.md
├── shared/               # Shared Go modules
│   ├── config/
│   ├── health/
│   ├── http-response/
│   ├── logger/
│   └── middlewares/
├── data-service/         # Foundation layer
│   ├── README.md
│   ├── Makefile
│   └── docker/
├── session-service/      # Authentication
│   ├── README.md
│   ├── Makefile
│   └── docker/
├── gateway-service/      # API Gateway
│   ├── Makefile
│   └── docker/
└── ui-service/           # Frontend
    ├── Makefile
    └── docker/
```

## Technology Stack

- **Backend**: Go 1.21+
- **Frontend**: HTML5, CSS3, JavaScript, Bootstrap 5
- **Database**: PostgreSQL 15
- **Containerization**: Docker & Docker Compose
- **Web Server**: Nginx (UI), Go net/http (services)

## Health Endpoints

All services expose health endpoints for monitoring:

| Service | Health Endpoint |
|---------|----------------|
| Gateway | `GET /api/v1/gateway/p/health` |
| Session | `GET /api/v1/sessions/p/health` |
| Data | `GET /api/v1/data/p/health` |
| UI | `GET /health` |

## Network

All services communicate through the `docker_barrest_network` Docker network.

## Development

### Go Workspace

This project uses Go workspaces for managing multiple services. The workspace allows you to work on all services simultaneously:

```bash
# Work on any service from project root
go test ./invoice-service/...
go build ./gateway-service/

# Or navigate to service directory
cd invoice-service && go test ./...
```

See [Workspace Development Guide](docs/6.%20workspace.md) for detailed instructions.

### Service Development

For service-specific development instructions, see each service's README:

- [Data Service](data-service/README.md) - Database setup, migrations
- [Session Service](session-service/README.md) - Authentication, JWT configuration
