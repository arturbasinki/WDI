---
type: Aggregate
title: Plot
description: "Aggregate of a plotted function: formula, interval, requested vs applied sample count, continuous segments, asymptote positions and the applied Y clamp."
tags: [aggregate, z06]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Aggregate. The rendered-ready description of one sampled function: what was plotted, over
what interval, with what sampling, split into continuous segments, annotated with asymptotes,
and clamped at a stated threshold.

# Schema

- `formula` — the expression as supplied (echoed back so the client can show what was drawn),
- `xmin`, `xmax` — the interval,
- `samplesRequested` — what the client asked for,
- `samples` — what was actually used after clamping to 501…20001,
- `segments` — list of point sequences, each continuous and never crossing an asymptote,
- `asymptotes` — x-positions of detected vertical asymptotes,
- `yclamp` — the ±threshold applied to outlier ordinates.

# Rules

- Sampling places points at equal intervals across the closed interval, endpoints included,
  so `samples` is the number of points per full curve.
- `samples` and `samplesRequested` are both returned: a client that asked for 5 000 000 points
  must be told it got 20 001.
- Asymptotes closer than two sample steps are merged, so a wide asymptotic neighbourhood
  yields one dashed line, not a cluster.
- Fewer than ten finite samples ⇒ an empty plot with a default clamp, not a crash.
- The interval is rejected when the lower bound is not strictly below the upper — written so
  that a NaN bound is rejected too.

# Related

[Function plotting](/domains/function-plotting.md) ·
[Plot segment](/model/plot-segment.md) · [Plot clamping](/model/plot-clamping.md) ·
[No silent truncation](/rules/no-silent-truncation.md) ·
[Plot a function](/usecases/plot-function.md)

# Citations

- `internal/z06/expr.go` — `Plot`, `PlotFun`: sample grid, percentile thresholds,
  segmentation, dedup, sample clamping, the `!(xmin < xmax)` rejection.
- `internal/z06/clamp_test.go` — branches reaching the frame, requested vs applied samples.
