#!/bin/sh
set -eu

: "${POSTGRES_URL:?Ошибка: POSTGRES_URL не задан (пример: postgres://user:pass@host:port/dbname?sslmode=disable)}"
MAX_WAIT_SECONDS="${MAX_WAIT_SECONDS:-60}"
WAIT_INTERVAL_SECONDS="${WAIT_INTERVAL_SECONDS:-2}"

# --- Функции логирования ---
# Формат: TIMESTAMP [LEVEL] [PID] Сообщение
timestamp() {
  date +"%Y-%m-%dT%H:%M:%S%z"
}

log_info() {
  printf '%s [INFO]  [%s] %s\n' "$(timestamp)" "$$" "$*"
}

log_warn() {
  printf '%s [WARN]  [%s] %s\n' "$(timestamp)" "$$" "$*"
}

log_error() {
  printf '%s [ERROR] [%s] %s\n' "$(timestamp)" "$$" "$*" >&2
}

# --- Парсинг URL и формирование ADMIN_URL ---
# Отсекаем схему (до "://") и оставшуюся часть
scheme=${POSTGRES_URL%%://*}
url_no_scheme=${POSTGRES_URL#*://}
credentials_and_host=${url_no_scheme%%/*}
rest_of_path="/${url_no_scheme#*/}"

db_path=${rest_of_path%%\?*}       # например "/mydb"
DB_NAME=${db_path#/}               # убираем ведущий "/"
if [ -z "$DB_NAME" ]; then
  log_error "Не удалось извлечь имя базы данных из POSTGRES_URL."
  exit 1
fi

# Сохраняем строку параметров (если есть)
query_string=""
case "$rest_of_path" in
  *\?*)
    query_string="?${rest_of_path#*\?}"
    ;;
esac

ADMIN_URL="${scheme}://${credentials_and_host}/postgres${query_string}"

# Для схемы pgx5 преобразуем в postgres:// для миграций
if [ "$scheme" = "pgx5" ]; then
  MIGRATE_URL=`echo "$POSTGRES_URL" | sed 's/^pgx5:\/\//postgres:\/\//g'`
else
  MIGRATE_URL="$POSTGRES_URL"
fi

# --- Ожидание готовности PostgreSQL ---
elapsed=0
while true; do
  if psql "$ADMIN_URL" -c '\q' >/dev/null 2>&1; then
    break
  fi

  if [ "$elapsed" -ge "$MAX_WAIT_SECONDS" ]; then
    log_error "Не удалось подключиться к PostgreSQL за ${MAX_WAIT_SECONDS} секунд."
    exit 1
  fi

  log_info "Ожидание PostgreSQL (${elapsed}/${MAX_WAIT_SECONDS})..."
  sleep "$WAIT_INTERVAL_SECONDS"
  elapsed=$((elapsed + WAIT_INTERVAL_SECONDS))
done

# --- Проверка существования целевой БД и создание при необходимости ---
exists=$(psql "$ADMIN_URL" -Atqc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}';")
if [ -z "$exists" ]; then
  log_info "База \"${DB_NAME}\" не найдена. Создаём."
  if ! psql "$ADMIN_URL" -c "CREATE DATABASE \"${DB_NAME}\";"; then
    log_error "Не удалось создать базу \"${DB_NAME}\"."
    exit 1
  fi
else
  log_info "База \"${DB_NAME}\" уже существует."
fi

log_info "Запуск user_service."
exec /usr/local/bin/catalog_service
