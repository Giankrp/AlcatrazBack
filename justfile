# Justfile for AlcatrazBack

binary_name := "alcatraz-server"
main_package := "./cmd/server/main.go"

# List available tasks
default:
    @just --list

# Compile the executable
build:
    go build -o {{ binary_name }} {{ main_package }}

# Run the server directly
run:
    go run {{ main_package }}

# Run all project tests
test:
    go test ./... -v

# Clean generated binaries
clean:
    rm -f {{ binary_name }}

# Download and tidy module dependencies
tidy:
    go mod tidy

# Run security-specific tests
test-security:
    go test ./internal/security/... -v

# Start services with Docker Compose
up:
    docker compose up -d

# Stop Docker Compose services
down:
    docker compose down

# Show container logs
logs:
    docker compose logs -f

# Clear all table data (WARNING: Destructive!)
reset-db:
    @echo "⚠️ WARNING: ALL database data will be deleted in 3 seconds..."
    @sleep 3
    docker compose exec -T postgres psql -U postgres -d alcatraz -c "DO \$\$ DECLARE _table_name record; BEGIN FOR _table_name IN SELECT quote_ident(schemaname) || '.' || quote_ident(tablename) AS tbl FROM pg_tables WHERE schemaname = 'public' LOOP EXECUTE 'TRUNCATE TABLE ' || _table_name.tbl || ' CASCADE;'; END LOOP; END \$\$;"
    @echo "✅ Database cleared successfully."

# Run psql shell
sql:
    docker compose exec postgres bash
