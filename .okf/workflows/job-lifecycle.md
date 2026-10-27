---
type: Workflow
title: Job lifecycle workflow
description: "The multi-step process of a background computation: validate, acquire a slot, register, accept, compute with progress, terminate, release, expire."
tags: [workflow, jobs]
timestamp: "2026-10-27T16:32:44Z"
---

# Scope

The end-to-end life of one background computation, from submission to eviction.

# Steps

1. **Submit.** Client posts parameters for one of the three long capabilities.
2. **Validate** synchronously. Out of bounds → 400, no job, no slot taken.
3. **Acquire a slot.** Both slots busy → 429, no job created.
4. **Register.** A job is created in state `running` with an opaque id; the slot is already held.
5. **Accept.** The client receives 202 with the job id. The job's cancellation handle is stored
   under the job's lock.
6. **Compute.** A worker runs the solver with (a) a cancellation signal **not** derived from the
   HTTP request, and (b) a progress callback writing snapshots under the lock.
7. **Report.** The worker records the outcome: result → `done`, error → `error`, panic →
   `error` (never a dead process), cancellation observed → `cancelled`. Status and end time are
   written together.
8. **Release the slot** after the terminal status is visible.
9. **Poll.** The client reads snapshots: `running` + progress, then a terminal state.
10. **Cancel** (optional, any time before step 7): 202 "cancelling"; the job notices at its
    granularity and finishes `cancelled` with partial state.
11. **Expire.** 15 minutes after ending, the job is swept on the next registry access; beyond
    200 retained jobs, the oldest finished are swept too. Polling an evicted job is a 404.

# Cross-cutting rules

[Concurrency cap](/rules/concurrency-cap.md) · [Job retention](/rules/job-retention.md) ·
[Cancellation semantics](/rules/cancellation-semantics.md) ·
[Progress cadence](/rules/progress-cadence.md) ·
[Job status monotonicity](/invariants/job-status-monotonic.md) ·
[Detached job context decision](/decisions/detached-job-context.md)

# Citations

- `cmd/wdi-server/main.go` — the job start helper, the registry (lazy sweep, semaphore),
  the shared GET/DELETE handlers, shutdown cancellation.
- `cmd/wdi-server/handlers_test.go` — the lifecycle asserted end to end with an injected clock.
- `cmd/wdi-server/web/app.js` — `startJob` and `poll` drive steps 5–10 from the client.
