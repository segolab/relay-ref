#!/usr/bin/env sh
set -eu

RELAY_BASE_URL="${RELAY_BASE_URL:-http://localhost:8429}"
RELAY_API_KEY="${RELAY_API_KEY:-dev-key}"

# Burst = 20, RPS = 10 → expect 429 after ~25 rapid requests
for i in $(seq 1 30); do
  curl -s -o /dev/null -w "%{http_code}\n" \
    -X POST "$RELAY_BASE_URL/v1/relays" \
    -H "X-API-Key: $RELAY_API_KEY" \
    -H "Content-Type: application/json" \
    -d @payload.json
done
