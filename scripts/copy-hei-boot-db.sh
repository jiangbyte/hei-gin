#!/usr/bin/env bash
set -euo pipefail

TARGET_DB="${1:-hei_gin}"
SOURCE_DB="hei_boot"

echo "==> Copy ${SOURCE_DB} -> ${TARGET_DB} via pg_dump (boot schema, no job migration)"

docker exec dev-postgres psql -U postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname='${TARGET_DB}' AND pid <> pg_backend_pid();" || true
docker exec dev-postgres psql -U postgres -c "DROP DATABASE IF EXISTS ${TARGET_DB};"
docker exec dev-postgres psql -U postgres -c "CREATE DATABASE ${TARGET_DB} OWNER postgres;"

docker exec dev-postgres pg_dump -U postgres -d "${SOURCE_DB}" --no-owner --no-acl \
  | docker exec -i dev-postgres psql -U postgres -d "${TARGET_DB}" -v ON_ERROR_STOP=1 -q

echo "==> Verify sys_job (boot columns)"
docker exec dev-postgres psql -U postgres -d "${TARGET_DB}" -c "SELECT id, name, handler, trigger_type FROM sys_job ORDER BY id;"

echo "Done: database ${TARGET_DB} ready for hei-gin (boot-aligned schema)"
