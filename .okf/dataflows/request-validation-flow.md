---
type: DataFlow
title: Request validation flow
description: "How a request travels from the browser form through intake, decoding, edge validation and interpretation to a capability, and where each rejection branch exits."
tags: [dataflow, api, validation]
timestamp: "2026-10-27T16:32:44Z"
---

# Path

Browser form → transport edge → per-capability validation → either a rejection with a reason,
or the capability → JSON response.

# Stages and their transformations

1. **Body intake** — the request body is bounded (4 MB); an over-long body never reaches parsing.
2. **Decoding** — JSON to a typed request; undecodable input is a 400 quoting the decode
   failure.
3. **Edge validation** — parameter ranges and text bounds, per capability
   ([Synchronous validation](/rules/synchronous-validation.md)); failure returns a
   [reason with the range](/contracts/error-response.md) and, for async capabilities,
   **no job and no slot**.
4. **Interpretation** — text is parsed into a domain object (a system, an expression, an
   equation, a grid); failure names the position.
5. **Computation** — synchronous capability answers, or the async path takes a slot and starts
   a [Job](/model/job.md).
6. **Response** — success payload with applied values echoed
   ([No silent truncation](/rules/no-silent-truncation.md)).

# Where the flow is enforced

Every branch is asserted at the transport layer with the real computation replaced by a
deterministic stub, so the flow is testable without waiting on minutes-long work
([Solver capability port](/contracts/solver-package-port.md)).

# Citations

- `cmd/wdi-server/main.go` — body limit, decode helper, per-endpoint validation, response
  helper.
- `cmd/wdi-server/handlers_test.go` — each rejection branch, with its message.
