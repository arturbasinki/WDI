---
type: BoundedContext
title: HTTP API context
description: "The bounded context serving the product: transport, edge validation, job policy, task catalog, static delivery, and test seams over the heavy computations."
tags: [bounded-context, hexagonal, adapter]
timestamp: "2026-10-27T16:32:44Z"
---

# Boundary

The single served process: it delivers the static frontend, exposes the JSON API, owns the
[Background computation](/domains/background-computation.md) policy, and enforces the input
bounds presented to learners.

# Responsibilities and their ownership

1. **Transport.** Routing, request-body size limit, JSON encode/decode, HTTP status mapping,
   Polish human-readable error messages.
2. **Validation at the edge.** Every parameter is validated synchronously, before a job exists
   and before it occupies a computation slot ([Synchronous validation](/rules/synchronous-validation.md)).
3. **Bounds that belong to the web layer.** z02's accepted limit of 1000…10¹² matches the UI
   selector; the solvers deliberately do not know these numbers.
4. **Job policy.** Concurrency cap, retention, cancellation, shutdown
   ([Background job policy](/rules/background-job-policy.md)).
5. **The task catalog** as the one writer of task identity
   ([Task catalog domain](/domains/task-catalog.md)).
6. **Test seams.** The three heavy computations are invoked through replaceable call points,
   so HTTP behaviour is tested deterministically without real computation.

# Deliberate omissions

No database, no cache layer, no authentication, no configuration file. The only runtime
setting is the listening port, taken from the environment. Graceful shutdown is required by
the container platform ([Graceful shutdown workflow](/workflows/graceful-shutdown.md)).

# Related

[HTTP API contract](/contracts/synchronous-solver-api.md) ·
[Background job API](/contracts/background-job-api.md) ·
[Request validation flow](/dataflows/request-validation-flow.md)

# Citations

- `cmd/wdi-server/main.go` — route table, body limit, job store, catalog, per-endpoint
  validation, shutdown wiring, `PORT` from environment.
- `cmd/wdi-server/api_test.go`, `handlers_test.go`, `z08_ui_test.go` — the HTTP behaviour
  contract, one test per acceptance criterion.
