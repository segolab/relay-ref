#!/usr/bin/env sh
set -eu

RELAY_BASE_URL="${RELAY_BASE_URL:-http://localhost:8429}"

curl -i \
  -X POST "$RELAY_BASE_URL/v1/relays" \
  -H "Content-Type: application/json" \
  -d @payload.json
