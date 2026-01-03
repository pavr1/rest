# 🔐 Session Service

Authentication and session management service for Bar-Restaurant application.

## Prerequisites

- Go 1.24+
- PostgreSQL running (via data-service)
- Sessions table created (via migration)

## Quick Start

### 1. Ensure Data Service is Running

```bash
cd ../data-service
make fresh           # Start and run migrations (or make start if already set up)
```

### 2. Start Session Service

```bash
cd ../session-service
make start           # Start container
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/sessions/p/login` | Staff login |
| `POST` | `/api/v1/sessions/p/validate` | Validate session |
| `POST` | `/api/v1/sessions/logout` | Logout |
| `GET` | `/api/v1/sessions/p/health` | Health check |

## Usage Examples

### Login

```bash
curl -X POST http://localhost:8087/api/v1/sessions/p/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

**Response:**
```json
{
  "code": 201,
  "message": "Login successful",
  "data": {
    "session_id": "abc123...",
    "token": "eyJhbG...",
    "staff": {
      "id": "uuid",
      "username": "admin",
      "first_name": "Admin",
      "last_name": "User",
      "role": "admin"
    }
  }
}
```

### Validate Session

```bash
curl -X POST http://localhost:8087/api/v1/sessions/p/validate \
  -H "Content-Type: application/json" \
  -d '{"session_id": "abc123..."}'
```

### Logout

```bash
curl -X POST http://localhost:8087/api/v1/sessions/logout \
  -H "Content-Type: application/json" \
  -d '{"session_id": "abc123..."}'
```

### Health Check

```bash
curl http://localhost:8087/api/v1/sessions/p/health
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make start` | Start session-service container |
| `make stop` | Stop container |
| `make restart` | Restart container |
| `make logs` | View container logs |
| `make status` | Show container status |
| `make clean` | Remove container and images |
| `make test` | Run unit tests |
| `make help` | Show all commands |

## Configuration

Settings are loaded from the data-service `settings` table:

| Key | Default | Description |
|-----|---------|-------------|
| `JWT_SECRET` | (generated) | JWT signing key |
| `JWT_EXPIRATION_TIME` | `24h` | Token expiration |
| `SERVER_HOST` | `0.0.0.0` | Service host |
| `SERVER_PORT` | `8087` | Service port |

## Staff Roles

- `admin` - Full access
- `manager` - Management access
- `waiter` - Order management
- `bartender` - Bar operations
- `chef` - Kitchen operations
- `dj_karaoke_operator` - Karaoke management

## Default Admin User

Created by data-service init script:
- **Username:** `admin`
- **Password:** `admin123`

## Directory Structure

```
session-service/
├── entities/sessions/
│   ├── handlers/
│   │   ├── db_handler.go    # Database operations
│   │   ├── http_handler.go  # HTTP endpoints
│   │   └── jwt_handler.go   # JWT token handling
│   ├── models/
│   │   └── models.go        # Data structures
│   └── sql/
│       ├── queries.go       # SQL query loader
│       └── scripts/*.sql    # SQL queries
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
├── scripts/
│   ├── build.sh
│   └── start.sh
├── main.go
├── main_http_handler.go
├── Makefile
└── go.mod
```

## Troubleshooting

### "Failed to load configuration"
Ensure data-service is running and accessible.

### "Failed to connect to database"
Check database connection settings in data-service settings table.

### "Invalid username or password"
Verify staff credentials. Default: `admin` / `admin123`
