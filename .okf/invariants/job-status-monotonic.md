---
type: Invariant
title: Job status monotonicity
description: "A job's status is monotonic: running then exactly one terminal state, with the end timestamp written atomically with it."
tags: [invariant, jobs, concurrency]
timestamp: "2026-10-27T16:32:44Z"
---

# Statement

A job moves from `running` to exactly one terminal status (`done`, `error`, `cancelled`),
never twice, and the end timestamp is written in the same atomic step as the terminal status.
A job in a terminal state cannot be cancelled and cannot be re-ended.

# Why it must hold

Clients poll status and act on it. If a finished job could be marked cancelled, the interface
would show a cancelled result that was in fact delivered. If the end timestamp could be observed
before the status, retention could evict a job that has not finished.

# Observed failures this prevents

- DELETE on an already finished job returned "cancelling" — a promise that nothing was going to
  happen. It is now a conflict with a human-readable reason.
- A data race on the cancellation field (written by the accepting request, read by DELETE and
  by shutdown) was detected by the race detector; the field is now lock-protected.
- A search that was cancelled failed to terminate: the worker pool's input channel was not
  closed on cancellation, so workers waited forever and the job never reached a terminal state.

# Related

[Job](/model/job.md) · [Job lifecycle workflow](/workflows/job-lifecycle.md) ·
[Background job policy](/rules/background-job-policy.md) · [Race freedom](/quality/race-freedom.md)

# Citations

- `cmd/wdi-server/main.go` — `finish` (status + end time under one lock), `markCancelled`
  (only while running), `tryCancel` (false when finished), panic → `error`.
- `internal/z02/armstrong.go` — the worker channel is always closed, so cancellation terminates.
- `cmd/wdi-server/handlers_test.go` — 202 → cancelled → 409, and the finished-job conflict.
