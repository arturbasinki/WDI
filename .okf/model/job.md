---
type: Aggregate
title: Job
description: "Aggregate for one background computation: opaque id, monotonic status lifecycle, progress/result snapshots, cancellation handle and end time, all lock-protected."
tags: [aggregate, jobs, concurrency]
timestamp: "2026-10-27T16:32:44Z"
---

# Type and identity

Aggregate root of [Background computation](/domains/background-computation.md). Identity: a
generated job id, opaque to clients, unique within the process, handed out at acceptance and
used as the resource address for polling and cancellation.

# Lifecycle (status)

`running` → exactly one of `done` | `error` | `cancelled`. Terminal states are reached once
and are never revisited.

# Schema

- status, optional error message, latest progress snapshot, final result,
- a cancellation handle,
- the moment it ended — which drives retention.

**Every field, including the cancellation handle and the end timestamp, is protected by the
job's own lock.** A POST writes the handle while DELETE and shutdown read it; without the lock
this is a classic data race, and it was detected by the race detector.

# Rules the aggregate enforces

- Ending a job sets status and end time **together**, so retention can never observe an
  "ended" job without a date.
- Marking a job cancelled applies only while it is still running; a job that finished in the
  meantime keeps its real state.
- Cancelling is attempted, not forced: an already-finished job reports "nothing to interrupt",
  which the transport renders as a conflict.
- A crash inside the computation becomes status `error`, not a dead process: the worker
  goroutine is outside the transport's own recovery, so an unrecovered crash would kill the
  whole server.
- The result exposed to clients is a snapshot built under the lock, never a live reference.

# Related

[Job lifecycle workflow](/workflows/job-lifecycle.md) ·
[Search progress](/model/search-progress.md) ·
[Background job policy](/rules/background-job-policy.md) ·
[Detached job context decision](/decisions/detached-job-context.md)

# Citations

- `cmd/wdi-server/main.go` — the job type with its mutex, `setCancel`, `tryCancel`,
  `snapshot`, `setProgress`, `finish`, `markCancelled`, and the panic-recovering worker.
- `cmd/wdi-server/handlers_test.go` — race-sensitive lifecycle, 202→cancelled→409 sequence.
