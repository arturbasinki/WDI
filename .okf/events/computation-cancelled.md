---
type: DomainEvent
title: Computation cancelled
description: "Domain event: a running computation was interrupted by client, navigation or shutdown, carrying its partial state and freeing its slot."
tags: [domain-event, jobs, cancellation]
timestamp: "2026-10-27T16:32:44Z"
---

# Event

A running computation was interrupted, and the partial state at the moment of interruption is
known. Named in the past tense: the computation **was cancelled**.

# Sources of the cancellation

1. A client asked to cancel the job ([Cancellation semantics](/rules/cancellation-semantics.md)).
2. The client left the view that owned the job.
3. The process is shutting down ([Graceful shutdown workflow](/workflows/graceful-shutdown.md)).

# Emitted when

The capability observed the cancellation and returned; the job's status becomes `cancelled`
only if it was still running — a job that finished in the meantime keeps its real state.

# Payload carried

The partial state: numbers found so far, the keys tried, the measurements completed. The
capability returns the partial state **together with** the cancellation, and the transport does
not present it as a result
([Cancellation granularity](/rules/cancellation-granularity.md)).

# Who reacts

- **Polling client** — sees `status: "cancelled"`, stops polling.
- **Slot release** — the freed slot lets a queued learner start.
- **Shutdown sequence** — cancelling everything live is what allows the process to exit inside
  the grace window.

# Termination obligation

Cancelling must unblock the machinery that produces the event: a worker pool waiting on an
input channel that is never closed means the event never fires and the job hangs forever.

# Citations

- `cmd/wdi-server/main.go` — `markCancelled`, `tryCancel`, `cancelAll`, the worker's
  `ctx.Err()` branch.
- `internal/z02/armstrong.go` — the always-closed channel; `internal/z03/sieve.go` — the
  in-measurement checks.
- `cmd/wdi-server/handlers_test.go`, `internal/z02/cancel_test.go`, `internal/z03/cancel_test.go`.
