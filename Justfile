config := "env.yaml"

default:
    @just --list

# ── Development ──────────────────────────────────────────────────────────────

# Run in watch mode
run:
    CONFIG_PATH={{config}} wgo run ./cmd/api

# Run in watch mode with race detection
run-race:
    CONFIG_PATH={{config}} wgo run -race ./cmd/api

# Run with test config
run-test:
    CONFIG_PATH=env.test.yaml go run -race ./cmd/api

# Build binary
build:
    go build -o dist/api ./cmd/api

# Start built binary
start:
    CONFIG_PATH={{config}} ./dist/api


# ── Testing ───────────────────────────────────────────────────────────────────

# Run unit tests with race detection
test:
    go run gotest.tools/gotestsum --format testdox -- -race ./internal/... ./cmd/... ./config/...

# Run e2e tests against a running API (start it with `just run-test` in another shell)
e2e base_url="http://127.0.0.1:8080":
    E2E_BASE_URL={{base_url}} go run gotest.tools/gotestsum --format testdox -- -race ./e2e/...


# ── Linting ───────────────────────────────────────────────────────────────────

# Format all Go code
fmt:
    go fmt ./...

# Apply go fix
fix:
    go fix ./...


# ── Postgres / sqlc ───────────────────────────────────────────────────────────

# Generate sqlc code from SQL queries
sqlc-gen:
    sqlc generate -f internal/infrastructure/storage/postgres/sqlc.yaml

# Run database migrations
# Commands: up (default), down, version, force
# Options:  steps=N (number of steps), version=N (for force)
migrate command="up" steps="-1":
    go run cmd/migrate/main.go -command={{command}} -steps={{steps}}
