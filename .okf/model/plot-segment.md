---
type: ValueObject
title: Plot segment
description: "A maximal continuous stretch of a plotted curve, broken at non-finite samples and asymptotic jumps so the renderer never draws across an asymptote."
tags: [value-object, z06]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Value object. One maximal stretch of the curve that is continuous on the drawing surface:
an ordered sequence of points, with no asymptote and no non-finite value inside it.

# Why the curve is split at all

A single polyline through an asymptote draws a vertical spike joining +∞ to −∞ — the classic
plotting artefact the exercise calls out. Returning segments hands the decision to the
renderer: each segment is one path, each asymptote is one dashed line.

# Where a segment breaks

- A sample whose ordinate is infinite or NaN ends the segment; if the neighbouring sample is
  finite, an asymptote sits between them.
- A sign change between neighbours whose jump exceeds 4·IQR of |y| ends the segment and marks
  an asymptote.
- Entering or leaving the outlier zone does **not** break a segment any more: outlier values
  are clamped, so the curve stays continuous right up to the frame
  ([Plot clamping](/model/plot-clamping.md)).

# Related

[Plot](/model/plot.md) · [Numeric asymptote detection](/decisions/numeric-asymptote-detection.md)

# Citations

- `internal/z06/expr.go` — `Point`, `Plot.Segments`, the `breakAt` scan and the
  segment-building loop.
- `internal/z06/clamp_test.go` — `1/x` on [−5,5] yields two branches reaching ±8.36.
