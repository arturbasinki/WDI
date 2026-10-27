---
type: ValueObject
title: Elimination step
description: "A single narrated elimination operation (description + full matrix snapshot) used to animate the reduction, bounded by a cell budget."
tags: [value-object, z01, animation]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Value object. One recorded operation on the matrix during elimination: a human-readable
description plus a full copy of the matrix at that moment.

# Purpose

The narration is a deliverable, not a log: the interface replays the steps as an animation of
the reduction, so the learner sees *how* the answer was obtained.

# Recorded operations

Row swap (leading-element selection), row normalisation (divide the leading row by its pivot),
and row subtraction (eliminate a column from every other row — full Gauss-Jordan, both above
and below). The initial augmented matrix is recorded as step zero.

# Cost and its bound

Each step stores a complete n×(n+1) matrix copy and a full elimination performs ≈ n² steps, so
the narration grows as n⁴. A budget expressed in **matrix cells** (400 000) bounds it; when
exhausted, registration stops and `stepsTruncated` is set. The computed answer is unaffected.

# Related

[Linear system](/model/linear-system.md) ·
[Animation step budget](/rules/animation-step-budget.md) ·
[Resource boundedness](/quality/resource-boundedness.md)

# Citations

- `internal/z01/solver.go` — `Step`, `snapshot`, the `record` closure with the cell budget,
  `StepsTruncated`.
- `cmd/wdi-server/web/app.js` — the z01 view replays `steps` as an animation.
- `internal/z01/stepsbudget_test.go` — truncation flag on a large system.
