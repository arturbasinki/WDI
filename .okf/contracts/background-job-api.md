---
type: Contract
title: Background job API
description: "The shared submit/poll/cancel contract for the three long-running capabilities, with the jobId, progress snapshot and cancellation shapes."
tags: [contract, api, async, jobs]
timestamp: "2026-10-27T16:32:44Z"
---

# Port

The three long-running capabilities share one asynchronous shape: submit, poll, cancel.

# Endpoints

| Capability | Submit | Poll | Cancel |
|---|---|---|---|
| Narcissistic-number search | `POST /api/z02/search` (`limit`) | `GET /api/z02/search/{id}` | `DELETE /api/z02/search/{id}` |
| Sieve benchmark | `POST /api/z03/benchmark` (`ns`) | `GET /api/z03/benchmark/{id}` | `DELETE /api/z03/benchmark/{id}` |
| Cipher break | `POST /api/z05/break` (`cipher`, `maxLen`, `capKeys`) | `GET /api/z05/break/{id}` | `DELETE /api/z05/break/{id}` |
| Verified number list | `GET /api/z02/known` | — | — |

`GET` and `DELETE` for all three are one shared behaviour, not three implementations.

# Success payloads

- Submit → 202 `{ "jobId": "…" }`
- Poll → 200 `{ jobId, status, progress?, result?, error? }` where status ∈
  `running | done | error | cancelled`
- Cancel running → 200 `{ "status": "cancelling" }`

# Status contract

Full table in [Background job policy](/rules/background-job-policy.md). The status codes are
part of the domain, asserted one per acceptance criterion: 400 validation before a job exists,
429 at the concurrency cap, 404 for an evicted id, 409 for cancelling a finished job.

# Client obligations

Poll on an interval; stop on a terminal state, on failure (with a visible banner), and on
leaving the view; cancel the job when its view is left
([View resource disposal](/rules/view-resource-disposal.md)).

# Related

[Job lifecycle workflow](/workflows/job-lifecycle.md) · [Job](/model/job.md) ·
[Progress flow](/dataflows/progress-flow.md)

# Citations

- `cmd/wdi-server/main.go` — the job start helper, the shared GET/DELETE handlers, the three
  submit handlers.
- `cmd/wdi-server/handlers_test.go` — the status contract.
- `cmd/wdi-server/web/app.js` — `startJob`, `poll`, cancel-on-leave.
