# Root-level variables
config := "env.yaml"

# List all available commands
default:
    @just --list

# Run in watch mode
run:
    CONFIG_PATH={{config}} wgo run cmd/aircraft_api/main.go

# Run in watch mode with -race
run-race:
    CONFIG_PATH={{config}} wgo run -race cmd/aircraft_api/main.go

# Build
build:
    CONFIG_PATH={{config}} go build -o dist/aircraft_api cmd/aircraft_api/main.go

# Start built binary
start:
    CONFIG_PATH={{config}} ./dist/aircraft_api

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


# --- Postgres Ent ---

# Generate ent code from schemas
ent-gen:
    go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert --target ./internal/infrastructure/storage/postgres/ent ./internal/infrastructure/storage/postgres/schema
