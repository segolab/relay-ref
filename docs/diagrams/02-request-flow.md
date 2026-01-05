# Diagram 02 — Request flow

Rate limiting & idempotency (`POST /v1/relays`)

## Purpose of this diagram

* Show **exact request flow** for creating a relay
* Make **rate limiting** and **idempotency** explicit
* Clarify **what happens on retries**
* Explain **why ordering of concerns matters**

This diagram focuses on **ingress behavior only**.

```mermaid
sequenceDiagram
    participant Client
    participant API as Relay API
    participant RL as Rate Limiter
    participant IDEM as Idempotency Store
    participant STORE as Relay Store

    Client->>API: POST /v1/relays
    API->>RL: Check rate limit (apiKey, route)

    alt Rate limit exceeded
        RL-->>API: Reject
        API-->>Client: 429 Too Many Requests\nRetry-After
    else Allowed
        RL-->>API: Allow

        API->>IDEM: Lookup Idempotency-Key
        alt Existing relay
            IDEM-->>API: Existing relay
            API-->>Client: 201 Created (same relay)
        else New request
            IDEM-->>API: Not found
            API->>STORE: Create relay
            STORE-->>API: Relay created
            API-->>Client: 201 Created
        end
    end
```

### How to read this diagram

1. Rate limiting comes first
    * Applied before idempotency
    * Protects the system from overload
    * Ensures fairness per API key
    * Failed requests still consume rate-limit budget
1. Idempotency is checked only after allowance
    * Prevents duplicate relays on retries
    * Same idempotency key returns the same relay
    * Conflicting payloads result in an error (not shown here)
1. Storage is the final step
    * Relay is created only once
    * Listing and inspection rely on this stored state
