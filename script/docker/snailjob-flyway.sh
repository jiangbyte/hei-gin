#!/usr/bin/env bash
# One-shot Flyway migrate for snail_job on existing local Postgres (dev-postgres).
# Does NOT recreate or reconfigure the Postgres container.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
FLYWAY_IMAGE="${FLYWAY_IMAGE:-swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/flyway/flyway:11-azure-mongo}"
PG_HOST="${SNAIL_JOB_PG_HOST:-host.docker.internal}"
PG_PORT="${SNAIL_JOB_PG_PORT:-5432}"
JDBC_URL="jdbc:postgresql://${PG_HOST}:${PG_PORT}/snail_job"
JDBC_USER="${SNAIL_JOB_DB_USER:-admin}"
JDBC_PASSWORD="${SNAIL_JOB_DB_PASSWORD:-123456}"

docker run --rm \
  --entrypoint flyway \
  --add-host=host.docker.internal:host-gateway \
  -v "${ROOT}/script/sql/postgres/snailjob:/flyway/sql:ro" \
  "${FLYWAY_IMAGE}" \
  -url="${JDBC_URL}" \
  -user="${JDBC_USER}" \
  -password="${JDBC_PASSWORD}" \
  -locations="filesystem:/flyway/sql" \
  -baselineOnMigrate=true \
  migrate
