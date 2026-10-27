---
type: BusinessRule
title: Synchronous validation
description: "Out-of-bounds parameters are rejected at submission time with a client error, before any job exists and before a computation slot is taken."
tags: [business-rule, api, validation]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

**Whenever** a request carries parameters outside their declared bounds, **then** reject it
synchronously with a client-error response — before a background job is created and before it
occupies a computation slot.

# Bounds enforced at the edge

| Capability | Parameter | Accepted range |
|---|---|---|
| z02 | `limit` | 1000 … 10¹² (0 means "use the default") |
| z03 | each `n` | 1000 … 100 000 000 |
| z05 | `maxLen` | 1 … 9 |
| z05 | `capKeys` | 1 … 100 000 000 |
| z05 | `mode` | `enc` \| `dec` |
| z05 | `cipher` | non-empty after trimming |
| z06 | `expr` | ≤ 2000 characters |
| z06 | interval | lower strictly below upper (NaN included in the rejection) |
| z07 | `limit` | ≥ 0, above the ceiling it is clamped, not rejected |
| z01 | parsed system | see [size invariant](/invariants/linear-system-size.md) |
| all | request body | ≤ 4 MB |

# Why this is a business rule, not validation hygiene

Accepting a bad parameter and failing later inside a job produces a job that exists, occupies
one of the two computation slots, runs for a while, and then reports an error the client could
have been told about immediately. The learner sees a spinner for a typo. Synchronous rejection
is the observable contract, and one test per bound pins it.

# Related

[Background job policy](/rules/background-job-policy.md) ·
[Layer ownership of bounds](/rules/layer-ownership-of-bounds.md) ·
[HTTP API context](/contexts/http-api.md) · [Error response](/contracts/error-response.md)

# Citations

- `cmd/wdi-server/main.go` — per-endpoint validation blocks preceding the job start helper;
  the body size limit.
- `cmd/wdi-server/handlers_test.go` — "400 at POST, no job created" for each bound.
- `.verify/gate.sh` — a live probe asserting an out-of-range limit yields 400.
