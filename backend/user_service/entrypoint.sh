#!/bin/sh
set -eu

: "${POSTGRES_URL:?Error: POSTGRES_URL is not set (e.g.: postgres://user:pass@host:port/dbname?sslmode=disable)}"

# Parse POSTGRES_URL into prefix (everything before last "/") and database name
PREFIX="${POSTGRES_URL%\/*}"
DB_WITH_PARAMS="${POSTGRES_URL#${PREFIX}/}"
DB_NAME="${DB_WITH_PARAMS%%\?*}"
# Admin URL (to connect to the default "postgres" database)
ADMIN_URL="${PREFIX}/postgres${POSTGRES_URL#*${DB_NAME}}"

timestamp() { date +"%Y-%m-%dT%H:%M:%S%z"; }
log_info()  { printf '\033[1;34m%s [INFO]  [%s]\033[0m %s\n' "$(timestamp)" "$$" "$*"; }
log_error() { printf '\033[1;31m%s [ERROR] [%s]\033[0m %s\n' "$(timestamp)" "$$" "$*" >&2; }

# 1) Wait until PostgreSQL at the admin URL becomes available
elapsed=0
MAX_WAIT_SECONDS="${MAX_WAIT_SECONDS:-60}"
WAIT_INTERVAL_SECONDS="${WAIT_INTERVAL_SECONDS:-2}"

log_info "Waiting for PostgreSQL to become available at $ADMIN_URL..."
while ! psql "$ADMIN_URL" -c '\q' >/dev/null 2>&1; do
  if [ "$elapsed" -ge "$MAX_WAIT_SECONDS" ]; then
    log_error "PostgreSQL did not become available within ${MAX_WAIT_SECONDS} seconds"
    exit 1
  fi
  log_info "Still waiting for PostgreSQL (${elapsed}/${MAX_WAIT_SECONDS})..."
  sleep "$WAIT_INTERVAL_SECONDS"
  elapsed=$((elapsed + WAIT_INTERVAL_SECONDS))
done

# 2) For local development: drop and recreate the database
log_info "Dropping database \"$DB_NAME\" (if it exists)..."
if ! psql "$ADMIN_URL" -c "DROP DATABASE IF EXISTS \"${DB_NAME}\";"; then
  log_error "Failed to execute DROP DATABASE \"${DB_NAME}\""
  exit 1
fi

log_info "Creating database \"$DB_NAME\"..."
if ! psql "$ADMIN_URL" -c "CREATE DATABASE \"${DB_NAME}\";"; then
  log_error "Failed to execute CREATE DATABASE \"${DB_NAME}\""
  exit 1
fi

# 3) Run golang-migrate from a clean slate
log_info "Running golang-migrate for database $DB_NAME..."
if ! migrate -path /migrations -database "$POSTGRES_URL" up; then
  log_error "Failed to apply migrations via golang-migrate"
  exit 1
fi
log_info "Migrations applied successfully"

# 4) Start the service
log_info "Launching user_service"
exec /usr/local/bin/user_service
