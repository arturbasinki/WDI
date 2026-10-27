---
type: QualityAttribute
title: Race freedom
description: "Race freedom: the full suite runs under a data-race detector in the gate; job state, progress producers, registry, slots and test seams are lock- and join-disciplined."
tags: [quality-attribute, concurrency, testing, iso25010]
timestamp: "2026-10-27T16:32:44Z"
---

# Attribute class (ISO/IEC 25010)

**Maintainability** + **Reliability** — the concurrent parts of the product are verified
data-race free, on every change.

# Target

The full test suite runs with the race detector as a gate step, so a data race is a build
failure, not a future incident.

# Concurrent surface and its discipline

- **Job state** — every field, including the cancellation handle and the end timestamp, is
  guarded by the job's own lock; status and end time are written in one critical section
  ([Job](/model/job.md)).
- **Progress producers** — the ticker is stopped and joined before the final result is
  assembled; progress payloads are copies, sorted outside the critical section
  ([Progress cadence](/rules/progress-cadence.md)).
- **Registry** — one lock for map and creation order; a job's own lock is taken briefly per
  entry, never held while iterating the map ([Job retention](/rules/job-retention.md)).
- **Slot accounting** — the slot is released after the terminal status is visible
  ([Concurrency cap](/rules/concurrency-cap.md)).
- **Worker pools** — input channel always closed on cancellation
  ([Cancellation granularity](/rules/cancellation-granularity.md)).
- **Test harness** — job goroutines are drained before a test's stubs are restored, so a stale
  goroutine cannot read a stub after its test ended.
- **Client side** — per-route disposers release intervals and handlers; route sequencing
  discards a stale render ([View resource disposal](/rules/view-resource-disposal.md)).

# Races found and closed by this discipline

The progress ticker racing the final sort; the cancellation handle written by POST and read by
DELETE without a lock; a search hanging forever after cancellation because its channel was
never closed.

# Citations

- `.github/workflows/ci.yml`, `.verify/gate.sh` — the race-enabled suite as a gate step.
- `cmd/wdi-server/main.go` — the mutex-guarded job, the registry sweep, the semaphore.
- `internal/z02/armstrong.go` — ticker join, always-closed channel.
- `cmd/wdi-server/handlers_test.go` — the drain-before-restore ordering, the concurrency probe.
