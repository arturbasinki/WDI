---
type: ValueObject
title: Plot clamping
description: "The IQR-derived ±Y threshold applied to outlier samples, applied by keeping them at the threshold so branches reach the frame of the plot."
tags: [value-object, z06, honesty]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Value object. The ± threshold applied to outlier ordinates of a plot, and the value that
results from applying it.

# Derivation

From the finite samples' |y| distribution: `YClamp = median + 15 · IQR`, where IQR is the
interquartile range; when IQR degenerates (below 1e-12) it is replaced by
`max(1, median)` so a constant function still has a usable threshold.

# The rule that makes it real

Clamping **keeps** outlier samples with their ordinate set to ±YClamp. Removing them instead —
the earlier behaviour — made a branch end in mid-air below the threshold: for `1/x` on [−5,5]
the branches stopped around |y| ≈ 41, halfway to the frame. With real clamping, the maximum
drawn |y| equals `YClamp` (8.36 for that case), so the branch visibly runs into the asymptote.

# Honesty

`YClamp` is returned in the plot so the client's axis and the reader's reading match the data.

# Related

[Plot](/model/plot.md) · [Plot segment](/model/plot-segment.md) ·
[Numeric asymptote detection decision](/decisions/numeric-asymptote-detection.md)

# Citations

- `internal/z06/expr.go` — percentile/quantile helper, `clampThresh`, the `clamped` array,
  `Plot.YClamp`.
- `internal/z06/clamp_test.go` — max |y| equals the clamp.
