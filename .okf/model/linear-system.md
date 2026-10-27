---
type: Aggregate
title: Linear system
description: "Aggregate of an augmented N×(N+1) matrix: square, exactly-supplied, finite coefficients, with rank/pivots/RREF and a parametric solution derived from it."
tags: [aggregate, z01]
timestamp: "2026-10-27T16:32:44Z"
---

# Type and identity

Aggregate root of [Linear systems](/domains/linear-systems.md). Identity: the parsed
coefficient matrix as supplied (N rows of N coefficients plus a free term); there is no
separate id — the same input text always denotes the same system.

# Schema

- `n` — number of equations and unknowns (square system).
- rows — N coefficient rows, each with the free term appended (an augmented matrix).
- Derived, once solved: rank, pivot columns, free columns, reduced row echelon form,
  the solution (values, or parametric expressions in the free parameters `t1…tk`).

# Invariants

- N ≥ 1 and bounded above (see [Linear system size bound](/invariants/linear-system-size.md)).
- Exactly N·(N+1) numeric tokens must be supplied — too few and too many are both errors.
- Every coefficient must be a finite number; NaN and infinities are rejected at parse time,
  because solving them yields a "unique" solution made of nonsense.
- Zero-test uses a fixed tolerance appropriate to the magnitudes of the exercise data
  (|v| < 1e-9 counts as zero).

# Parsing tolerance

Tokens may be separated by spaces, tabs, newlines or commas, and a sign separated from its
digits ("- 1 2 4 6") is merged with the following token, because that is how the exercise's
example file is written.

# Consistency boundary

The recorded elimination steps belong to the same aggregate as the answer: they narrate the
same reduction. Recording is budgeted, and truncating the narration never invalidates the
answer ([Animation step budget](/rules/animation-step-budget.md)).

# Citations

- `internal/z01/solver.go` — `Parse`, `Solve`, the augmented matrix, `Result` fields
  (`N`, `Status`, `X`, `RREF`, `FreeVars`, `SolutionLines`, `Steps`, `StepsTruncated`, `Rank`),
  finiteness check, token merging.
