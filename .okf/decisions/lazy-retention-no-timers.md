---
type: Decision
title: "Lazy job retention, no timers"
description: "Retain finished jobs for 15 minutes and at most 200, sweeping lazily on registry access with an injectable clock, instead of running background timers."
tags: [decision, adr, jobs, memory]
timestamp: "2026-10-27T16:32:44Z"
---

# Context

Job snapshots were never removed. A benchmark job holds rows for every structure at every
range; measured at roughly 0.4 MB per job, a few hundred finished jobs was ~200 MB of resident
memory with no visible cause and no restart pressure to reveal it.

# Decision

Keep finished jobs for 15 minutes, cap the registry at 200, and **sweep lazily on access** —
on creating a job and on looking one up — with no background timers and no housekeeping thread.
Expiry is a side effect of use.

# Alternatives considered

A sweeper goroutine on a ticker (rejected: a second lifecycle to reason about, and a clock the
tests must outwit); sweeping on a request-count heuristic (rejected: not time-exact); a fixed
per-job size cap alone (rejected: does not bound count).

# Consequences

**Positive.** Deterministic and testable: a test injects a clock, advances it past the window,
and asserts the next lookup is a 404 — without sleeping. No timers means nothing to stop at
shutdown. The sweep preserves creation order, so "oldest finished first" is exact.

**Negative.** Expiry is deferred until someone touches the registry; on an idle server memory is
held longer than 15 minutes. Accepted: the process is single-purpose and idle retention costs
nothing observable.

**Related.** [Job retention](/rules/job-retention.md) · [Job expired](/events/job-expired.md) ·
[Job](/model/job.md) · [Resource boundedness](/quality/resource-boundedness.md)

# Citations

- `cmd/wdi-server/main.go` — the retention window, the count ceiling, `sweepLocked`,
  `removeLocked`, the injectable `now` on the registry.
- `cmd/wdi-server/handlers_test.go` — the fake clock, TTL expiry, the capacity ceiling.
- `RAPORT-JAKOSC-KODU.md` §5 P1 item 1 — the leak this decision closed.
