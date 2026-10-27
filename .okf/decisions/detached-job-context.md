---
type: Decision
title: Job context detached from the request
description: "A job's cancellation is detached from the HTTP request that created it, so long computations survive connection churn and end only on deliberate cancellation."
tags: [decision, adr, jobs, cancellation]
timestamp: "2026-10-27T16:32:44Z"
---

# Context

The first implementation ran every long computation under the request's own cancellation
context. A client that closed the POST — a page refresh, a navigation, a proxy closing an idle
connection — killed a search that had been running for minutes, and the learner lost the work
with no way to see partial results.

# Decision

A job's cancellation signal is **detached** from the HTTP request. It is created when the job
starts and is cancelled only by:
- an explicit DELETE on the job,
- the client's cancel-on-leave, issued deliberately,
- process shutdown.

The POST answers immediately with a job id; the computation continues independently of the
connection that created it.

# Consequences

**Positive.** Long work survives connection churn. Cancellation becomes an intentional,
observable act with defined semantics for each job state
([Cancellation semantics](/rules/cancellation-semantics.md)). The learner can navigate away and
come back while the job lives (until retention).

**Negative.** Work continues after a client disappears, consuming one of the two computation
slots — which is precisely why the concurrency cap, the retention policy, and cancel-on-leave
are required alongside this decision.

**Confirmed by.** All three job endpoints; the cancellation lifecycle tests.

# Related

[Job](/model/job.md) · [Job lifecycle workflow](/workflows/job-lifecycle.md) ·
[Concurrency cap](/rules/concurrency-cap.md) ·
[View resource disposal](/rules/view-resource-disposal.md)

# Citations

- `cmd/wdi-server/main.go` — the detached context created in the job start helper (not derived
  from the request), the cancel handle stored under the job lock.
- `cmd/wdi-server/handlers_test.go` — the cancellation lifecycle after the POST completes.
- `cmd/wdi-server/web/app.js` — `startJob`'s keepalive DELETE on leaving the view.
