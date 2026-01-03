#!/bin/bash

# Reset the database (WARNING: Deletes all data!)

echo "⚠️  WARNING: This will delete ALL database data!"
read -p "Are you sure you want to continue? (y/N): " confirm

if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
    echo "❌ Reset cancelled."
    exit 0
fi

cd "$(dirname "$0")/../docker"

echo "🗑️  Stopping containers and removing volumes..."
docker-compose down -v

echo "🔄 Starting fresh containers..."
docker-compose up -d

echo "⏳ Waiting for database to be ready..."
sleep 5

until docker exec barrest_postgres pg_isready -U postgres -d barrest_db > /dev/null 2>&1; do
    echo "Waiting for PostgreSQL..."
    sleep 1
done

echo "✅ Database reset complete!"
