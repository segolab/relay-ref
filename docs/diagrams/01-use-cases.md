# Diagram 01 — Use-cases (Relay API)

## Purpose of this diagram

* Explain **who** uses the Relay API
* Clarify **what a "relay"** is in abstract terms
* Show **why endpoints are source–target agnostic**
* Establish the mental model used by all later diagrams

This is **not** about implementation, infra, or delivery.


```mermaid
flowchart LR
    Source["Event Source<br/>(Client / Service / System)"]
    RelayAPI["Relay API<br/>(Record relay intent)"]
    Target["Destination<br/>(Webhook / Service / System)"]

    Source -->|"POST /v1/relays<br/>Create relay intent"| RelayAPI
    Source -->|"GET /v1/relays<br/>List relays"| RelayAPI
    Source -->|"GET /v1/relays/{id}<br/>Inspect relay"| RelayAPI

    RelayAPI -.->|"Delivery observation<br/>(test / metrics / simulation)"| Target
```

### How to read this diagram

* **Source**  
Any client or system capable of submitting events:
    * backend services
    * batch jobs
    * scripts
    * external systems
* **Relay API**  
Accepts, validates, rate-limits, and stores relay requests.  
It does **not** execute delivery in the baseline.
* **Target**  
The intended destination of the relay:
    * webhook endpoint
    * downstream service
    * external system  
Targets are **described**, not contacted.
* **Dashed arrow**  
Indicates a conceptual future step, not implemented behavior.
