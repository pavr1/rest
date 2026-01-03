#!/bin/bash

# Start the data service containers

echo "🍺 Starting Bar-Restaurant Data Service containers..."

cd "$(dirname "$0")/../docker"

docker-compose up -d

echo "⏳ Waiting for database to be ready..."
sleep 3

# Check if postgres is ready
until docker exec barrest_postgres pg_isready -U postgres -d barrest_db > /dev/null 2>&1; do
    echo "Waiting for PostgreSQL..."
    sleep 1
done

echo "✅ Data Service containers started successfully!"
echo ""
echo "📝 Services available:"
echo "   PostgreSQL: localhost:5432"
echo "   PgAdmin:    http://localhost:8080"
echo "   Data API:   http://localhost:8086"
