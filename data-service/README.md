# Data Service

Database service for the Bar-Restaurant application. Provides PostgreSQL database with PgAdmin for management.

## Prerequisites

- Docker & Docker Compose
- Make

## Quick Start

```bash
# First time setup
chmod +x scripts/*.sh
make fresh
```

## Commands

| Command | Description |
|---------|-------------|
| `make fresh` | Clean install (deletes all data) |
| `make start` | Start containers (keeps data) |
| `make stop` | Stop containers |
| `make restart` | Restart containers |
| `make logs` | View logs |
| `make connect` | Connect to PostgreSQL CLI |
| `make status` | Show container status |
| `make clean` | Remove containers and volumes |
| `make migrate` | Apply pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-status` | Show applied migrations |
| `make test` | Run unit tests |

## Connection Info

### PostgreSQL
- **Host:** localhost
- **Port:** 5432
- **Database:** barrest_db
- **User:** postgres
- **Password:** postgres123

```
postgresql://postgres:postgres123@localhost:5432/barrest_db
```

### PgAdmin
- **URL:** http://localhost:8080
- **Email:** admin@barrest.com
- **Password:** admin123

To connect in PgAdmin:
- Open http://localhost:8080
- Login: admin@barrest.com / admin123
- Right-click "Servers" → "Register" → "Server"
- General tab:
- Name: Bar-Restaurant
- Connection tab:
- Host: postgres (not localhost or 127.0.0.1)
- Port: 5432
- Database: barrest_db
- Username: postgres
- Password: postgres123
- Click "Save"

## Database Schema

The database includes 25 tables for the bar-restaurant application:

**Core:** Tables, Customers, Orders, Order Items, Menu Items, Menu Categories, Stock Items, Stock Item Categories, Menu Item Stock Items, Suppliers, Purchase Invoices, Invoice Details, Existences, Customer Favorites

**Business Logic:** Payments, Customer Invoices, Staff, Request Notifications, Karaoke Song Requests, Karaoke Song Library, Loyalty Points Transactions, Table Sessions

**Advanced:** Reservations, Promotions, Reviews

## Migrations

Place migration files in `docker/init/migrations/`:

```
001_description.up.sql    # Apply
001_description.down.sql  # Rollback
```

Then run:
```bash
make migrate
```

## Troubleshooting

### Port already in use
```bash
make clean
make start
```

### Database won't start
```bash
make fresh  # Complete reset
```

### Check logs
```bash
make logs s=postgres
```
