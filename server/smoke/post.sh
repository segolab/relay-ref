#!/usr/bin/env sh
set -eu

RELAY_BASE_URL="${RELAY_BASE_URL:-http://localhost:8429}"
RELAY_API_KEY="${RELAY_API_KEY:-dev-key}"

curl -i \
  -X POST "$RELAY_BASE_URL/v1/relays" \
  -H "X-API-Key: $RELAY_API_KEY" \
  -H "Content-Type: application/json" \
  -d @payload.json
