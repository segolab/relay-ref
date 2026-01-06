# Smoke tests

These scripts perform manual, server-agnostic HTTP sanity checks.

They assume:
- a relay-ref server is running
- environment variables are set

They test:
- [basic POST](./post.sh)
- [idempotency](./post-idempotent.sh)
- [unauthorized access](./unauthorized.sh)
- [rate limiting behavior](./rate-limit.sh)

These are not load or performance tests.
