# Challenge: Delivery Ingestion

## Overview

The baseline Relay API records and observes **delivery intent** only.
Actual delivery or execution is intentionally out of scope.

This challenge replaces **delivery observation** with **real ingestion**
while preserving the existing API contract and client behavior.

The goal is to evolve the system without breaking:
- OpenAPI contract
- endpoint semantics
- rate limiting behavior
- client expectations

---

## Baseline behavior

In the baseline implementation:

- POST /v1/relays records a relay intent
- Relays are observable via GET endpoints
- No external calls are performed
- Delivery state is inferred, not executed

This provides testability and observability without side effects.

---

## Challenge goal

Introduce real delivery to the configured destination
(e.g. webhook URL) while keeping the API surface unchanged.

The Relay API must remain:
- source-agnostic
- target-agnostic
- idempotent
- rate-limited at ingress

---

## Constraints

- API endpoints and schemas must not change
- Delivery must not block request handling
- Failures must be observable
- Rate limiting semantics must remain intact
- The baseline must still work without delivery enabled

---

## Suggested difficulty levels

### Level 1 — Simulated ingestion (easy → medium)

- Introduce an internal delivery executor
- Transition relay state from `queued` to `delivered`
- No external I/O
- Useful for testing async behavior and observability

### Level 2 — Real webhook delivery (medium)

- Perform HTTP calls to destination URLs
- Add timeouts and basic retry logic
- Handle transient vs permanent failures
- Observe delivery outcomes via logs and metrics

### Level 3 — Production-grade ingestion (hard)

- Add persistent queueing
- Implement retry with backoff
- Introduce dead-letter handling
- Apply per-destination rate limiting
- Reason about delivery guarantees

---

## Out of scope

- Changes to client libraries
- Changes to OpenAPI contract
- Exactly-once delivery guarantees
- Provider-specific integrations

---

## Teaching focus

This challenge explores:

- evolving systems under a stable contract
- async execution and background workers
- observability-driven development
- operational complexity vs baseline simplicity
- trade-offs in reliability and throughput
