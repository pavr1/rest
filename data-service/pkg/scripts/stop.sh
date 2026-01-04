#!/bin/bash

# Stop the data service containers

echo "🛑 Stopping Bar-Restaurant Data Service containers..."

cd "$(dirname "$0")/../docker"

docker-compose down

echo "✅ Data Service containers stopped!"
