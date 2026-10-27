---
type: Policy
title: Cancellation semantics
description: "Cancellation semantics by job state: running cancels with an accepted acknowledgement, finished conflicts, unknown is not found; the client also cancels on leaving a view."
tags: [policy, jobs, api]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

**Whenever** a client asks to stop a job, **then**:
- the job is running → cancel it and acknowledge with "accepted" and status `cancelling`;
- the job is already in a terminal state → conflict, with a message saying there is nothing to
  interrupt;
- the job id is unknown (never existed, or already expired from retention) → not found.

The same three-way handling is shared by every cancellable capability; it is not reimplemented
per endpoint.

# Cancellation is cooperative, not immediate

Cancelling signals; the computation notices at its chosen granularity and returns partial state
with the cancellation. The acknowledgement therefore says "cancelling", not "cancelled": the
job reaches `cancelled` on the next poll.

# A cancellation the client does not issue

Leaving the page. The client cancels a running job when its view is left, so the learner's
navigation frees the computation slot. Process shutdown cancels every live job
([Graceful shutdown workflow](/workflows/graceful-shutdown.md)).

# Why the conflict case exists

Returning "cancelling" for a finished job was observed in practice: the client waited for a
state change that could never happen. A conflict with a reason is the honest answer.

# Related

[Job](/model/job.md) · [Job status monotonicity](/invariants/job-status-monotonic.md) ·
[Cancellation granularity](/rules/cancellation-granularity.md) ·
[Background job API](/contracts/background-job-api.md)

# Citations

- `cmd/wdi-server/main.go` — the shared job GET/DELETE handlers, the try-cancel guard, the
  shutdown cancellation of all jobs.
- `cmd/wdi-server/handlers_test.go` — active → 202 → `cancelled` → 409; unknown → 404.
- `cmd/wdi-server/web/app.js` — cancel-on-leave with a keepalive request.
