---
type: DomainEvent
title: Job expired
description: "Domain event: a finished job passed its retention window (or the count ceiling evicted it) and polling it now yields not-found."
tags: [domain-event, jobs, retention]
timestamp: "2026-10-27T16:32:44Z"
---

# Event

A finished job left the retention window and is no longer addressable. Named in the past tense:
the job **expired**.

# Emitted when

The lazy sweep runs — on creating a job or looking one up — and removes a job that ended more
than 15 minutes ago, or one of the oldest finished jobs when more than 200 are held
([Job retention](/rules/job-retention.md)). A running job never expires.

# Who reacts

- **Polling client** — the next poll is a 404 with "unknown job — expired, or never existed".
  A client that was still polling a finished job past its window is told the truth rather than
  being served a fabricated state.
- **Memory** — the snapshot is released; this is the mechanism that prevents unbounded growth
  from accumulated job results.

# Deliberate design choices behind the event

- No timers: expiry is a consequence of access, so the process runs no housekeeping goroutines
  and expiry is deterministic under an injected clock
  ([Lazy retention decision](/decisions/lazy-retention-no-timers.md)).
- Expiry of a *finished* job is invisible to a running one: the sweep never touches live work.
- Job state is ephemeral by design; nothing is persisted, so expiry is total.

# Citations

- `cmd/wdi-server/main.go` — `sweepLocked` (TTL pass then oldest-finished pass), `removeLocked`
  preserving creation order, the injectable clock.
- `cmd/wdi-server/handlers_test.go` — fake clock advanced past the window ⇒ 404; capacity
  eviction test.
