---
type: Policy
title: Cancellation granularity
description: Cancellation is observed inside long loops at a near-zero-cost granularity (once per 2^20 iterations) so a slow measurement is interrupted in a fraction of a second.
tags: [policy, cancellation, performance, z03]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

A long computation must notice cancellation **inside** its work, at a granularity whose cost is
measurably zero — not only between top-level phases.

# Granularity in force

Inside the sieve loops, the cancellation signal is tested once every 2²⁰ iterations, using a
full-bit mask so the test is a single comparison amortised over a million iterations. The check
sits in **every** inner loop, including the flag-initialisation pass: at n = 10⁸, initialising
the table alone costs ~0.1 s normally and ~7 s under the race detector, so checking only during
crossing-out would make cancellation wait out the initialisation.

# Measured effect

Cancelling a benchmark went from ~7 s to ~0.04 s. Without in-measurement checks, cancelling a
map-based sieve at n = 10⁸ (4…15 s per measurement) means waiting minutes for the whole
benchmark.

# Rule for the interrupted work

A measurement interrupted mid-way is **not** reported as a row: partial rows would read as
"this structure is slow". The benchmark returns the rows completed so far together with the
cancellation. Likewise the number search returns partial progress plus the cancellation, and
never reports success.

# Termination of the workers

Cancellation must also unblock the machinery: the pool's input channel is always closed on the
cancellation path, otherwise workers wait on it forever and the job never finishes — the
observed "search hangs forever after cancellation" defect.

# Related

[Prime sieving](/domains/prime-sieving.md) · [Search progress](/model/search-progress.md) ·
[Job status monotonicity](/invariants/job-status-monotonic.md) ·
[Cancellation-inside-measurement decision](/decisions/cancellation-inside-measurement.md)

# Citations

- `internal/z03/sieve.go` — `ctxTick = 1 << 20`, the `canceled` mask test in every loop,
  the partial-rows-plus-error return.
- `internal/z03/cancel_test.go` — cancellation interrupts a single measurement.
- `internal/z02/armstrong.go` — `defer close(jobs)` on the producer, partial state with
  `ctx.Err()`.
- `internal/z02/cancel_test.go` — cancelled search terminates with the partial result.
