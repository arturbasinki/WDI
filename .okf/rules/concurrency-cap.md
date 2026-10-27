---
type: Policy
title: Computation concurrency cap
description: At most two computations run concurrently; the third submission is refused with a too-many-requests response and no job is created.
tags: [policy, jobs, resource]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

**Whenever** a long computation is requested and the maximum number of concurrent computations
is already running, **then** refuse the request with "too many requests" naming the cap and
telling the client to retry — do not queue it, and do not create the job.

# The cap

Two concurrent computations. The slot is taken **before** the job is registered, and released
**after** the job's terminal status has been written, so a refusal never interleaves with a job
that is still starting, and a finished job's slot is not freed while its status is not yet visible.

# Rationale

Each heavy sweep is core-saturating and memory-hungry (a narcissistic-number search to a large
limit is a full-core, minutes-long run). Six accepted searches at once exhaust CPU and memory
and degrade every other learner. Refusing is honest and immediate; queueing would look like a
job that hangs.

# Interaction with cancellation

The slot is released by the worker's own completion path, including the crash path, so a failed
computation cannot leak the slot. Cancellation ends the computation, which ends the worker,
which frees the slot — that is why "cancel to make room" is a legitimate client strategy.

# Related

[Background job policy](/rules/background-job-policy.md) ·
[Job lifecycle workflow](/workflows/job-lifecycle.md) ·
[Resource boundedness](/quality/resource-boundedness.md)

# Citations

- `cmd/wdi-server/main.go` — the counting-slot semaphore, `tryAcquire`/`release`, refusal in
  the job start helper, release in the deferred recovery block.
- `cmd/wdi-server/handlers_test.go` — two accepted then the third is refused; slot freed after
  completion.
