# Database Migrations

Place migration files here following this naming convention:

```
001_description.up.sql    # Apply migration
001_description.down.sql  # Rollback migration
```

## Example

```
001_add_order_notes.up.sql
001_add_order_notes.down.sql
002_add_customer_preferences.up.sql
002_add_customer_preferences.down.sql
```

## Commands

```bash
make migrate         # Apply all pending migrations
make migrate-down    # Rollback last migration
make migrate-status  # Show applied migrations
```

## Notes

- Migrations are tracked in the `schema_migrations` table
- Always create both `.up.sql` and `.down.sql` files
- Migrations run in alphabetical order (use numeric prefixes)

