#!/bin/bash

# Migration script for Bar-Restaurant database
# Usage: ./migrate.sh [up|down|status]

CONTAINER="barrest_postgres"
DB_NAME="barrest_db"
DB_USER="postgres"
MIGRATIONS_DIR="$(dirname "$0")/../docker/init/migrations"

# Ensure migrations directory exists
mkdir -p "$MIGRATIONS_DIR"

# Create schema_migrations table if not exists
docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -c "
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);" 2>/dev/null

case "$1" in
    up)
        echo "📈 Applying migrations..."
        for file in $(ls "$MIGRATIONS_DIR"/*.up.sql 2>/dev/null | sort); do
            version=$(basename "$file" .up.sql)
            
            # Check if already applied
            applied=$(docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -t -c \
                "SELECT COUNT(*) FROM schema_migrations WHERE version = '$version';" | tr -d ' ')
            
            if [ "$applied" = "0" ]; then
                echo "  Applying: $version"
                docker cp "$file" $CONTAINER:/tmp/migration.sql
                if docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -f /tmp/migration.sql; then
                    docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -c \
                        "INSERT INTO schema_migrations (version) VALUES ('$version');"
                    echo "  ✅ Applied: $version"
                else
                    echo "  ❌ Failed: $version"
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
