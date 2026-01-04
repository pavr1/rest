#!/bin/bash

# Connect to PostgreSQL database

echo "🔌 Connecting to Bar-Restaurant database..."

docker exec -it barrest_postgres psql -U postgres -d barrest_db
