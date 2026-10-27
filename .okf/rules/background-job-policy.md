---
type: Policy
title: Background job policy
description: "The complete status-code policy for background jobs: 202 accept, 400 validation, 429 capacity, 404 unknown, 202 cancelling, 409 finished, 200 snapshots."
tags: [policy, jobs, api, contract]
timestamp: "2026-10-27T16:32:44Z"
---

# The policy in one table

Every interaction with a background job maps to a status code and a human-readable Polish
message. This table is the contract; each row has a regression test written one-to-one with the
acceptance criterion.

| Situation | Response | Payload |
|---|---|---|
| Submission accepted, computation started | 202 Accepted | `{ "jobId": "…" }` |
| Parameters out of bounds | 400 | `{ "error": "<range stated, value stated>" }` — no job id |
| Both computation slots busy | 429 | `{ "error": "busy: N computations already running — try again shortly" }` |
| Poll a live job | 200 | `{ jobId, status: "running", progress: {…} }` |
| Poll a finished job | 200 | `{ jobId, status: "done", result: {…} }` |
| Poll a failed job | 200 | `{ jobId, status: "error", error: "…" }` |
| Poll a cancelled job | 200 | `{ jobId, status: "cancelled" }` |
| Poll an evicted/unknown job | 404 | `{ "error": "unknown job — expired or never existed" }` |
| Cancel a running job | 200 | `{ "status": "cancelling" }` |
| Cancel a finished job | 409 | `{ "error": "job already finished — nothing to interrupt" }` |
| Cancel an unknown job | 404 | as above |

# Invariants this table rests on

[Job status monotonicity](/invariants/job-status-monotonic.md),
[Synchronous validation](/rules/synchronous-validation.md),
[Concurrency cap](/rules/concurrency-cap.md),
[Job retention](/rules/job-retention.md),
[Cancellation semantics](/rules/cancellation-semantics.md).

# Shape note

Polling a job is a 200 even when the computation failed: the *job* answered correctly, and the
failure is the answer. Only the job's own outcome distinguishes `done` / `error` / `cancelled`.

# Related

[Background job API](/contracts/background-job-api.md) ·
[Error response contract](/contracts/error-response.md) ·
[Job lifecycle workflow](/workflows/job-lifecycle.md)

# Citations

- `cmd/wdi-server/main.go` — the job start helper, the shared GET/DELETE handlers, the
  JSON helpers for success and for error.
- `cmd/wdi-server/handlers_test.go` — one test per row, incl. the Polish conflict message.
