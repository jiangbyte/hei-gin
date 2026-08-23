#!/usr/bin/env bash
# Deprecated wrapper — use hei-boot/scripts/copy-hei-boot-db-mysql.sh (MySQL / dev-mysql).
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
exec bash "${ROOT}/hei-boot/scripts/copy-hei-boot-db-mysql.sh" "${1:-hei_gin}"
