---
type: Policy
title: Job retention policy
description: "Finished jobs are kept 15 minutes and at most 200 at a time, evicted lazily on access, oldest finished first, never a running job."
tags: [policy, jobs, memory]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

Finished jobs are retained for **15 minutes** after they end, and no more than **200** jobs are
held at all. Eviction removes first the jobs past their retention window, then, if still over
capacity, the oldest finished jobs in creation order. **A running job is never evicted.**

# Enforcement is lazy

Cleanup runs only when the registry is touched (on create and on lookup). There are no
background timers: the sweep is a side effect of use, which keeps the process free of
housekeeping goroutines and makes expiry deterministic to test.

# Why a ceiling at all

Job snapshots are the only state the process keeps, and they are not small — a benchmark job
holds rows, and a number-search job holds results; measured at roughly 0.4 MB per job.
Hundreds of accumulated jobs is a slow, unbounded memory leak with no visible cause.

# Observable consequences

- Polling an evicted job is a not-found with the message "expired, or never existed".
- Retention is tested with an injected clock, so no test sleeps for fifteen minutes.
- Capacity is asserted by counting retained jobs after pushing past the ceiling.

# Related

[Job](/model/job.md) · [Job lifecycle workflow](/workflows/job-lifecycle.md) ·
[Lazy retention decision](/decisions/lazy-retention-no-timers.md) ·
[Background job policy](/rules/background-job-policy.md)

# Citations

- `cmd/wdi-server/main.go` — the retention window, the count ceiling, the lazy sweep (TTL pass
  then oldest-finished pass), the injectable clock, `removeLocked` keeping creation order.
- `cmd/wdi-server/handlers_test.go` — fake clock advancing past TTL, capacity ceiling test.
