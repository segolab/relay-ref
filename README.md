# Relay Ref

Relay Ref is a reference implementation of a rate-limited relay
(message / webhook ingestion) API.

The project is intended as a reference architecture rather than a
feature-complete product. It focuses on how the same API contract is
implemented, tested, and operated across different ecosystems and
infrastructure choices.

---

## Goals

- Define a stable OpenAPI contract for a relay / message-style API
- Treat rate limiting as a first-class concern
- Compare Go and .NET ecosystem best practices
- Demonstrate rate-limit-aware clients
- Provide a foundation for testing, benchmarking, and IaC
- Keep the implementation small, readable, and reviewable

---

## Scope (v1 baseline)

- Small JSON payload relay (enqueue only)
- Per-API-key and per-route rate limiting
- Idempotent request handling
- Explicit 429 / Retry-After semantics
- Health and readiness endpoints

Actual delivery or execution of relays is intentionally out of scope
for the baseline.

---

## Repository layout

```text
api/        OpenAPI specification (source of truth)
server/     Server implementations (Go, .NET)
client/     Client implementations (Go, .NET)
tests/      Contract, integration, and load tests
infra/      Infrastructure as Code (Terraform first, others later)
docs/       Design notes and extension paths
```

---

## Reference-style development model

This repository follows a fixed, story-shaped commit history
(approximately 6–8 commits) that represents the conceptual layers of the
reference architecture:

- API contract
- Server baseline
- Client baseline
- Testing baseline
- Infrastructure baseline
- Documentation

Each commit is treated as a semantic chapter, not a chronological log.
Commits may be amended or extended as the reference evolves, but their
meaning is preserved.

We adhere to this adjustable baseline until the initial end-to-end lanes
are implemented, tested, and reviewed.

---

## Extensibility

The relay API is deliberately designed to support future complexity
without breaking the contract, such as:

- OCR / transcription job submission
- Metrics or log ingestion
- Git or CI event relaying
- Larger payloads or payload references

These extensions are considered out of scope for the baseline, but are
documented as future paths.

---

## Status

Early reference baseline.

History may be rewritten while the initial reference is being prepared.
