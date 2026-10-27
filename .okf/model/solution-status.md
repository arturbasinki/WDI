---
type: ValueObject
title: Solution status
description: "The three-way classification of a linear system (unique, contradictory, underdetermined) and the rendering rules that follow from each state."
tags: [value-object, z01]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Value object. The classification of a solved linear system — exactly one of three states,
compared by value, never mutated.

# States

| State | Meaning | How it is recognised |
|---|---|---|
| `unique` | exactly one solution | rank = N |
| `none` | contradictory, no solution | a row whose coefficients are all zero with a non-zero free term |
| `infinite` | underdetermined, infinitely many | rank < N and no contradictory row |

# Rendering rules that travel with the state

- `unique` — one line per unknown with its value; float noise is removed by rounding to nine
  decimals and a below-tolerance value is rendered as exactly 0.
- `none` — a single line stating that the system is contradictory.
- `infinite` — the free columns become parameters `t1…tk`; each pivot unknown is expressed as
  `free term − Σ (coefficient · parameter)`, so the sign of each parameter's contribution is
  the **negation** of the matrix coefficient. A coefficient of magnitude 1 is written as the
  bare parameter. This sign convention was once inverted, which is why it is pinned by a test.

# Related

[Linear system](/model/linear-system.md) · [Linear systems domain](/domains/linear-systems.md)

# Citations

- `internal/z01/solver.go` — the status constants, the contradictory-row scan, the parametric
  expression builder (`coef := -a`), the value rounding.
- `internal/z01/param_test.go` — the parametric sign convention.
