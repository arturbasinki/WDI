---
type: Policy
title: Animation step budget
description: "Recorded elimination steps are bounded by 400,000 matrix cells, with a truncation flag, so animation narration cannot exhaust memory."
tags: [policy, z01, memory, safety]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

The narration of an elimination is bounded by a budget of **400 000 matrix cells**. When the
budget is exhausted, step registration stops and `stepsTruncated` is set. The computed answer
is never affected.

# Rationale (this was the largest defect found)

Every recorded step stores a full n×(n+1) matrix copy, and a full Gauss-Jordan reduction
performs ≈ n² steps, so narration memory grows as n⁴. For n = 100 that is ~10 000 snapshots:
~1.1 GB of JSON response and ~4.7 GB resident memory — an out-of-memory kill of the process,
reachable by a single legitimate request. After the fix, the same request answers in tens of
milliseconds at ~1.18 MB.

# Shape of the rule

- The budget is expressed in cells, not in steps, so it bounds memory directly and identically
  for every n.
- Truncating narration is disclosed, not silent
  ([No silent truncation](/rules/no-silent-truncation.md)).
- The client's animation replays the steps it received; it is not told the animation is shorter
  than the derivation, which is why the flag exists.

# Related

[Elimination step](/model/elimination-step.md) ·
[Linear system](/model/linear-system.md) ·
[Resource boundedness](/quality/resource-boundedness.md) ·
[Linear system size](/invariants/linear-system-size.md)

# Citations

- `internal/z01/solver.go` — `maxStepCells = 400_000`, the `record` closure, `StepsTruncated`.
- `internal/z01/stepsbudget_test.go` — truncation on a large system.
- `.verify/gate.sh` — the n=100 response-size probe.
