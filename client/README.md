# Client Implementations

This directory contains client libraries for interacting with the Relay API.

Clients are expected to be rate-limit-aware and implement appropriate
retry and idempotency strategies.

Current and planned implementations:
- go/         Go client
- dotnet-cs/  .NET (C#) client
- dotnet-fs/  .NET (F#), planned
