#!/bin/bash

# Migration script for Bar-Restaurant database
# Usage: ./migrate.sh [up|down|status]

CONTAINER="barrest_postgres"
DB_NAME="barrest_db"
DB_USER="postgres"
# Get the directory where this script is located, then navigate to migrations directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATIONS_DIR="$SCRIPT_DIR/../../docker/init/migrations"

# Ensure migrations directory exists
mkdir -p "$MIGRATIONS_DIR"

# Check if container is running
if ! docker ps --format "table {{.Names}}" | grep -q "^${CONTAINER}$"; then
    echo "❌ Container $CONTAINER is not running"
    exit 1
fi

# Wait for database to be ready
echo "⏳ Waiting for database to be ready..."
max_attempts=30
attempt=1
while [ $attempt -le $max_attempts ]; do
    if docker exec $CONTAINER pg_isready -U $DB_USER -d $DB_NAME >/dev/null 2>&1; then
        echo "✅ Database is ready"
        break
    fi
    echo "   Attempt $attempt/$max_attempts: Database not ready yet..."
    sleep 1
    attempt=$((attempt + 1))
done

if [ $attempt -gt $max_attempts ]; then
    echo "❌ Database failed to become ready"
    exit 1
fi

# Create schema_migrations table if not exists
if ! docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -c "
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);" 2>/dev/null; then
    echo "❌ Failed to create schema_migrations table"
    exit 1
fi

case "$1" in
    up)
        echo "📈 Applying migrations..."
        echo "   Container: $CONTAINER"
        echo "   Database: $DB_NAME"
        echo "   User: $DB_USER"
        echo "   Migrations dir: $MIGRATIONS_DIR"

        migration_files=$(ls "$MIGRATIONS_DIR"/*.up.sql 2>/dev/null | sort)
        migration_count=$(echo "$migration_files" | wc -w)
        echo "   Found migration files: $migration_count"

        for file in $migration_files; do
            version=$(basename "$file" .up.sql)
            filename=$(basename "$file")

            echo "🔄 Processing migration: $filename (version: $version)"

            # Check if already applied
            echo "   Checking if already applied..."
            applied=$(docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -t -c \
                "SELECT COUNT(*) FROM schema_migrations WHERE version = '$version';" 2>/dev/null | tr -d ' \t\n\r')

            if [ "$applied" = "0" ] || [ -z "$applied" ]; then
                echo "   📄 Applying migration file: $filename"
                echo "   🔧 Executing SQL..."
                if docker cp "$file" $CONTAINER:/tmp/migration.sql && \
                   docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -f /tmp/migration.sql; then
                    echo "   💾 Recording migration in schema_migrations table..."
                    if docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -c \
                        "INSERT INTO schema_migrations (version) VALUES ('$version');" 2>/dev/null; then
                        echo "  ✅ Successfully applied: $version"
                    else
                        echo "  ❌ Failed to record migration: $version"
                        exit 1
                    fi
                else
                    echo "  ❌ Failed to apply migration SQL: $version"
                    exit 1
                fi
            else
                echo "  ⏭️  Skipping (already applied): $version"
            fi
        done
        echo "✅ Migrations complete!"
        ;;
    
    down)
        echo "📉 Rolling back last migration..."
        last_version=$(docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -t -c \
            "SELECT version FROM schema_migrations ORDER BY applied_at DESC LIMIT 1;" | tr -d ' ')
        
        if [ -z "$last_version" ]; then
            echo "No migrations to rollback."
            exit 0
        fi
        
        down_file="$MIGRATIONS_DIR/${last_version}.down.sql"
        if [ -f "$down_file" ]; then
            echo "  Rolling back: $last_version"
            docker cp "$down_file" $CONTAINER:/tmp/migration.sql
            if docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -f /tmp/migration.sql; then
                docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -c \
                    "DELETE FROM schema_migrations WHERE version = '$last_version';"
                echo "  ✅ Rolled back: $last_version"
            else
                echo "  ❌ Rollback failed: $last_version"
                exit 1
            fi
        else
            echo "  ❌ Down file not found: $down_file"
            exit 1
        fi
        ;;
    
    status)
        echo "📊 Migration status:"
        docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -c \
            "SELECT version, applied_at FROM schema_migrations ORDER BY applied_at;"
        ;;
    
    *)
        echo "Usage: $0 [up|down|status]"
        echo "  up     - Apply all pending migrations"
        echo "  down   - Rollback last migration"
        echo "  status - Show applied migrations"
        exit 1
        ;;
esac
