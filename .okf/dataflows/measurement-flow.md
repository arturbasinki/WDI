---
type: DataFlow
title: Benchmark measurement flow
description: "How benchmark data is produced and consumed: fixed structure order per range, collection-forced timings, rows, progress ticks, and log-scaled rendering."
tags: [dataflow, z03, benchmark]
timestamp: "2026-10-27T16:32:44Z"
---

# Path

Chosen ranges × three structures → sequential timed runs → rows → progress → log-scaled bars.

# Stages

1. **Ranges** — the submitted list of `n`, defaulting to the exercise ladder 10³…10⁸
   ([Benchmark sieve](/usecases/benchmark-sieve.md)).
2. **Ordering** — for each range, the three structures run one after another: slice, map,
   bitmap. Order is fixed so the columns of the table are comparable.
3. **Preparation** — collection is forced before each measurement, and collection remains
   enabled during it, so a row reflects a realistic cost rather than a tuned one.
4. **Timing** — wall-clock duration around the run, plus the prime count from that run.
5. **Row emission** — a completed measurement becomes one row and one progress tick
   ([Sieve benchmark row](/model/sieve-benchmark-row.md)); a measurement interrupted by
   cancellation emits **no row**
   ([Cancellation granularity](/rules/cancellation-granularity.md)).
6. **Presentation** — bars on a logarithmic scale, fastest structure at 0%, slowest at 100%;
   sub-millisecond values rendered as "< 1 ms".

# Why the presentation matters to the data

At 10⁸ the map costs seconds while the slice and bitmap cost tens of milliseconds. On a linear
scale every bar but the map's reads as zero, which inverts the exercise's conclusion.

# Reference outcome

n = 10⁷: slice ≈ 46 ms, bitmap ≈ 45 ms, map ≈ 4 s.

# Citations

- `internal/z03/sieve.go` — the measurement loop, ordering, collection forcing, row/progress.
- `cmd/wdi-server/web/app.js` — the logarithmic bar helpers and `fmtSec` in the z03 view.
- `internal/z03/sieve_test.go` — reference prime counts.
