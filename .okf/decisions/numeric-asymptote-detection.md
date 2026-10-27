---
type: Decision
title: Numeric asymptote detection and real clamping
description: "Detect vertical asymptotes numerically from the sample distribution (median + 15·IQR, 4·IQR jump, non-finite adjacency), segment the curve at them, and clamp outliers instead of dropping them."
tags: [decision, adr, z06, plotting]
timestamp: "2026-10-27T16:32:44Z"
---

# Context

The exercise requires plots of functions such as `1/x` and `tan(x)` to account for vertical
asymptotes, on a system that draws sampled points, not symbolic curves. Symbolic analysis was
off the table with the standard-library-only decision
([Single-binary decision](/decisions/stdlib-only-single-binary.md)).

# Decision

Detect asymptotes **numerically from the sample distribution**:
- threshold outliers by the distribution of |y| — median + 15·IQR, with a degenerate IQR
  replaced by `max(1, median)`;
- treat a sign change between neighbours whose jump exceeds 4·IQR as an asymptote;
- treat a non-finite sample adjacent to a finite one as an asymptote;
- merge asymptotes closer than two sample steps.

Then split the curve into **segments** that never cross an asymptote, and **clamp** outlier
ordinates to ±YClamp while keeping those samples.

# Alternatives and the trap

The obvious rule — drop samples beyond a threshold — was the earlier implementation and it is
visibly wrong: branches stopped in mid-air (for `1/x` on [−5,5] at |y| ≈ 41, half the way to
the frame). Clamping keeps the point at the threshold so the branch runs into the asymptote,
which is the behaviour the exercise's wording implies.

# Consequences

Asymptote positions are approximate to the sampling resolution; a dense oscillation near a
singularity can be missed at low sample counts; the thresholds are tuned constants rather than
theorems. In exchange, the plot never draws across a singularity and never truncates a branch,
and the applied clamp is disclosed
([Plot clamping](/model/plot-clamping.md)).

# Related

[Function plotting](/domains/function-plotting.md) · [Plot](/model/plot.md) ·
[Plot segment](/model/plot-segment.md) · [Plot a function](/usecases/plot-function.md)

# Citations

- `internal/z06/expr.go` — the quantile helper, `jumpThresh`, `clampThresh`, the `breakAt`
  scan, dedup, the `clamped` array.
- `internal/z06/clamp_test.go` — `1/x` on [−5,5]: clamp 8.36, max drawn |y| 8.36.
