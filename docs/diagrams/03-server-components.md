# Diagram 03 — Server components

## Purpose of this diagram

* Show what the server is responsible for
* Make separation of concerns explicit
* Explain where rate limiting and idempotency live
* Prepare ground for testing and observability
* Stay implementation-agnostic (Go / .NET)

This diagram is about structure, not flow.

```mermaid
flowchart TB
    Client[Client]

    subgraph Server["Relay Server"]
        HTTP["HTTP API<br/>Handlers"]
        AUTH["Auth & Request Context"]
        RL["Rate Limiter"]
        IDEM["Idempotency Store<br/>In memory"]
        STORE["Relay Store<br/>In memory"]
        OBS["Observation<br/>Logs and metrics"]
    end

    Client -->|HTTP requests| HTTP
    HTTP --> AUTH
    AUTH --> RL
    RL --> IDEM
    IDEM --> STORE

    HTTP --> OBS
    RL --> OBS
    IDEM --> OBS
    STORE --> OBS
```

### How to read this diagram

#### HTTP API

* Owns endpoint semantics
* Maps requests and responses
* Does not contain business state

#### Auth & Request Context

* Extracts API key
* Manages request ID / correlation
* Establishes request scope

#### Rate Limiter

* Enforces fairness and protection
* Applied before idempotency and storage
* Emits rate-limit signals for observation

#### Idempotency Store (baseline)

* In-memory, process-local
* Prevents duplicate relay creation
* TTL-bounded
* Explicitly **not persistent** in baseline

#### Relay Store (baseline)

* Stores relay intents
* In-memory, thread-safe
* Source of truth for GET endpoints

#### Observation

* Logs, metrics, counters
* Observes intent, not delivery
* Must degrade gracefully (no backend required)
