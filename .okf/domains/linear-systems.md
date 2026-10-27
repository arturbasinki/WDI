---
type: Domain
title: Linear systems (z01)
description: "Solve N×N systems of linear equations by pivoted Gauss-Jordan reduction, classifying unique, contradictory and underdetermined cases, and record the reduction steps."
tags: [domain, core, z01, numerics]
timestamp: "2026-10-27T16:32:44Z"
---

# Kind of subdomain

Core (the exercise itself). Solve a system of N linear equations in N unknowns and explain
the reduction.

# Problem space

Input: N, then N rows of N coefficients plus the free term. Output: the solution, classified
into exactly one of three cases — unique, contradictory (no solution), underdetermined
(infinitely many, expressed parametrically through the free unknowns). The exercise explicitly
requires the latter two cases to be recognised.

# Approach and its intent

Full Gauss-Jordan elimination with partial pivoting (largest magnitude in the column becomes
the leading element) yields reduced row echelon form, from which rank, pivot columns and free
columns follow directly. Classification is then mechanical: a row of all zeros with a non-zero
free term is a contradiction; rank equal to N is unique; anything else is parametric.

The elimination is recorded operation by operation (row swap, row normalisation, row subtraction)
so the interface can animate the reduction — the explanation is part of the deliverable, not a
by-product.

# Boundedness

The recording is the expensive part: every recorded step copies the whole matrix. A memory
budget on recorded cells stops registration and sets a truncation flag, leaving the computed
answer correct — see [Animation step budget](/rules/animation-step-budget.md) and
[Linear system size bound](/invariants/linear-system-size.md).

# Related

[Linear system aggregate](/model/linear-system.md) ·
[Solution status](/model/solution-status.md) ·
[Elimination step](/model/elimination-step.md) ·
[Solve a linear system](/usecases/solve-linear-system.md)

# Verified results

The exercise's own example yields x=18, y=−58, z=35; contradictory and parametric systems are
recognised. A 100×100 identity-like system is answered in tens of milliseconds with a
sub-megabyte response.

# Citations

- `internal/z01/solver.go` — parsing (incl. the "- 1" split-sign tolerance), elimination,
  classification, parametric solution lines, step recording and its cell budget.
- `internal/z01/solver_test.go`, `param_test.go`, `stepsbudget_test.go` — the three
  classifications, the parametric sign convention, the truncation flag.
