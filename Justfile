# Root-level variables
config := "env.yaml"

# List all available commands
default:
    @just --list

# Run in watch mode
run:
    CONFIG_PATH={{config}} wgo run cmd/api/main.go

# Run in watch mode with -race
run-race:
    CONFIG_PATH={{config}} wgo run -race cmd/api/main.go

# Build
build:
    CONFIG_PATH={{config}} go build -o dist/api cmd/api/main.go

# Start built binary
start:
    CONFIG_PATH={{config}} ./dist/api

# Run all tests
test:
    go test -v ./...

# Run tests with race detection
test-race:
    go test -v -race ./...

# Format all Go code
fmt:
    go fmt ./...

# Apply go fix
fix:
    go fix ./...


# --- Postgres sqlc ---

# Generate sqlc code from queries
sqlc-gen:
    sqlc generate -f internal/infrastructure/storage/postgres/sqlc.yaml

# Run database migrations (requires DATABASE_URL env var or -database-url flag)
# Commands:
#   up              - apply all pending migrations (default)
#   down            - rollback all migrations
#   down steps=N    - rollback N migrations
#   up   steps=N    - apply N migrations
#   version         - print current migration version and dirty state
#   force version=N - force-set version without running SQL (recover from dirty state)
migrate command="version" steps="-1":
    go run cmd/migrate/main.go -command={{command}} -steps={{steps}}
