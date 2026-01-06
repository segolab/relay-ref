#!/usr/bin/env sh
set -eu

RELAY_BASE_URL="${RELAY_BASE_URL:-http://localhost:8429}"
RELAY_API_KEY="${RELAY_API_KEY:-dev-key}"

IDEMP_KEY="smoke-idem-1"

curl -i \
  -X POST "$RELAY_BASE_URL/v1/relays" \
  -H "X-API-Key: $RELAY_API_KEY" \
  -H "Idempotency-Key: $IDEMP_KEY" \
  -H "Content-Type: application/json" \
  -d @payload.json

echo
echo "---- repeat ----"
echo

curl -i \
  -X POST "$RELAY_BASE_URL/v1/relays" \
  -H "X-API-Key: $RELAY_API_KEY" \
  -H "Idempotency-Key: $IDEMP_KEY" \
  -H "Content-Type: application/json" \
  -d @payload.json
