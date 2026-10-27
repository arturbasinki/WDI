---
type: Invariant
title: Linear system size and finiteness
description: "A linear system has 1..100 equations, exactly N*(N+1) finite numeric values, with NaN and infinities rejected before solving."
tags: [invariant, z01]
timestamp: "2026-10-27T16:32:44Z"
---

# Statement

For a supplied linear system: N ≥ 1, N ≤ 100, exactly N·(N+1) numeric values follow, and every
coefficient and free term is a **finite** number.

# Rationale

N is the size of the elimination (O(N³) time, O(N²) memory per recorded step). The upper bound
of 100 keeps a single request from exhausting memory once the narration is budgeted. Finiteness
must be checked at the boundary: the number parser accepts NaN and infinities, and solving a
matrix containing them yields status `unique` with a solution of NaN — a plausible-looking
answer that is meaningless.

# Enforcement

At parse time, before any arithmetic. Rejections name the offending row and value.

# Related

[Linear system](/model/linear-system.md) ·
[Animation step budget](/rules/animation-step-budget.md) ·
[Resource boundedness](/quality/resource-boundedness.md)

# Citations

- `internal/z01/solver.go` — `Parse`: the 1..100 check, the exact token-count checks (too few
  and too many), the `IsNaN`/`IsInf` rejection.
- `internal/z01/solver_test.go` — malformed input cases.
