---
type: UseCase
title: Search narcissistic numbers
description: "Start a cancellable background search for narcissistic numbers up to a validated bound, poll its progress, and read the found numbers."
tags: [use-case, z02, async]
timestamp: "2026-10-27T16:32:44Z"
---

# Actor · goal

Learner — find as many narcissistic/Armstrong numbers as possible up to a chosen bound, and
watch the search proceed.

# Main flow (background job)

1. Submit a search limit. `0` means the default of 10¹⁰.
2. Synchronous validation: 1000 … 10¹², matching the UI selector; outside it → 400 and no job
   ([Layer ownership of bounds](/rules/layer-ownership-of-bounds.md)).
3. Take a computation slot; if none is free → 429 ([Concurrency cap](/rules/concurrency-cap.md)).
4. Receive 202 with a job id.
5. Poll: progress carries `processed`, `limit`, and the numbers found so far, sorted, about
   ten times a second ([Progress cadence](/rules/progress-cadence.md)).
6. On completion, poll once more for `status: done` and the final list.

# Alternative flows

- **Cancel** → DELETE, 202 "cancelling", then `status: cancelled` with the numbers found so
  far; the search terminates promptly (the worker pool is unblocked, not left waiting).
- **Slot busy** → 429, retry shortly.
- **Leaving the page** → the client cancels the job, freeing the slot.
- **Server stopping** → all live jobs are cancelled
  ([Graceful shutdown workflow](/workflows/graceful-shutdown.md)).

# Companion read

The list of the 41 verified numbers to 10¹² is served directly, with the note that it was
verified by exhaustive search in this implementation; values ≤ 10⁹ are also covered by unit
tests ([Reproducibility](/quality/reproducibility.md)).

# Citations

- `internal/z02/armstrong.go` — `Search`, progress ticks, cancellation contract.
- `cmd/wdi-server/main.go` — the search endpoint with validation, the known-numbers endpoint
  carrying the verified list.
- `cmd/wdi-server/handlers_test.go` — 400/429/cancel paths; `internal/z02/*_test.go`.
