---
type: DomainEvent
title: Computation completed
description: "Domain event: a background computation finished successfully, its result became available, its retention window started and its computation slot was released."
tags: [domain-event, jobs]
timestamp: "2026-10-27T16:32:44Z"
---

# Event

A background computation finished on its own terms and its result is available.
Named in the past tense: the computation **completed**.

# Emitted when

The capability returned without error and the job was still running. The job's status becomes
`done`, the result is stored, and the end time is written in the same step
([Job status monotonicity](/invariants/job-status-monotonic.md)).

# Who reacts

- **Polling client** — the next snapshot carries `status: "done"` and the result; the poller
  stops.
- **Retention policy** — the end time starts the 15-minute window
  ([Job retention](/rules/job-retention.md)).
- **Slot release** — the computation slot is freed after the status is visible
  ([Concurrency cap](/rules/concurrency-cap.md)).

# Non-events (deliberately not this event)

A computation that was cancelled, and a computation that returned an error, are **not**
completions; they have their own outcomes. Publishing a completed result for a cancelled search
was a real defect (the number search reported `Done: true` mid-attack).

# Citations

- `cmd/wdi-server/main.go` — `finish` (result + status + end time under one lock), the worker's
  completion path, the slot release.
- `internal/z02/armstrong.go` — `Done: true` only on the joined, final state.
