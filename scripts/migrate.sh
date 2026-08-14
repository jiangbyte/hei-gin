#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CONFIG="${1:-$ROOT/config.yaml}"
cd "$ROOT"
go run ./cmd/migrate -config "$CONFIG" up
