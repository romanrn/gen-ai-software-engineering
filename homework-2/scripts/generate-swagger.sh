#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
SRC_DIR="$ROOT_DIR/src"
OUT_DIR="$ROOT_DIR/docs"
SWAG_BIN="$(go env GOPATH)/bin/swag"
SWAG_CMD="swag"

if command -v swag >/dev/null 2>&1; then
  SWAG_CMD="swag"
elif [[ -x "$SWAG_BIN" ]]; then
  SWAG_CMD="$SWAG_BIN"
else
  echo "Installing swag CLI..."
  (cd "$SRC_DIR" && go install github.com/swaggo/swag/cmd/swag@latest)
  SWAG_CMD="$SWAG_BIN"
fi

echo "Generating OpenAPI YAML from code annotations..."
(
  cd "$SRC_DIR" && \
  "$SWAG_CMD" init \
    -g main.go \
    -d cmd/api,internal/adapters/in/http,internal/domain \
    -o "$OUT_DIR" \
    --ot yaml
)

echo "Swagger generated at: $OUT_DIR/swagger.yaml"
