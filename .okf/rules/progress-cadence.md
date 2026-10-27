---
type: Policy
title: Progress cadence
description: "Progress is published at a fixed, low-cost cadence (10 Hz, per 500k keys, per measurement) with copies, a done flag reserved for finished runs, and joined producers."
tags: [policy, jobs, progress, concurrency]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

Progress is reported at a cadence that is meaningful to a watching human and negligible to the
computation — never per candidate.

# Cadences in force

| Computation | Cadence |
|---|---|
| Narcissistic-number search | ≈10 times a second, from a dedicated ticker |
| Cipher dictionary attack | every 500 000 keys tried (≈50 notifications/second at attack speed) |
| Sieve benchmark | once per completed measurement |

# Rules riding on the cadence

- `done: true` is emitted only for a genuinely finished run; a mid-attack notification carries
  `done: false`. An audit found an attack reporting `Done: true` while keys were still being
  tried — the client then stopped polling and displayed an incomplete answer.
- The progress payload is a **copy** and is sorted before it is sent, and the progress producer
  is stopped and joined **before** the final result is assembled. A final sort racing with the
  ticker's read of the same collection is a data race the race detector catches.
- A nil progress callback is legal everywhere: progress is optional, not a required collaborator.
- Progress is emitted at the producer's granularity; the transport stores the latest snapshot
  under a lock and hands out copies on poll ([Job](/model/job.md)).

# Related

[Search progress](/model/search-progress.md) · [Job](/model/job.md) ·
[Progress flow](/dataflows/progress-flow.md) · [Race freedom](/quality/race-freedom.md)

# Citations

- `internal/z02/armstrong.go` — the 100 ms ticker, the join before the final sort, sorted copies.
- `internal/z05/gronsfeld.go` — the 500 000-key `notify` cadence and the `Done` flags.
- `internal/z03/sieve.go` — one `onTick` per completed row.
- `cmd/wdi-server/main.go` — `setProgress` under the job lock, snapshot on read.
