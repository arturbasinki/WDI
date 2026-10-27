---
type: Aggregate
title: Cryptarithm puzzle
description: "Aggregate of a cryptarithm: the equation, its letters and no-zero letters, the search space size, the full solution count and the (possibly capped) solution list."
tags: [aggregate, z07]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Aggregate. A letter-for-digit puzzle: the equation as written, the letters it uses, the
letters forbidden to take zero, and the set of digit assignments that satisfy it.

# Schema

- `formula` — the normalised five-token equation `<lhs> <op> <rhs> = <result>`,
- `op` — one of `+ - * /`,
- `letters` — the distinct letters, ordered by code point,
- `spaceSize` — the size of the searched space, 10Pk for k letters,
- `count` — the **full** number of solutions,
- `truncated` — whether the returned list was cut,
- `solutions` — at most `limit` solutions, each a letter→digit mapping plus the equation
  rendered with digits,
- `timeMs` — search duration.

# Structural rules (the puzzle's identity)

- Exactly five whitespace-separated tokens with `=` in position four.
- Only uppercase A–Z in the three number tokens; each is at most 10 characters.
- At most ten distinct letters overall — eleven letters cannot be assigned ten digits.
- The first letter of any multi-digit number may not be zero.

# Search space pruning

Shape checks before searching: for `*`, the result's digit count must be
`len(lhs)+len(rhs)-1 … len(lhs)+len(rhs)`; for `/` the analogous constraint applies to the
dividend. A failing shape means zero solutions without searching. Multiplication overflow is
guarded (`B > MaxInt64 / A`) before comparison, so an overflow cannot masquerade as a match.

# Related

[Cryptarithm](/domains/cryptarithm.md) ·
[Distinct digit assignment](/invariants/distinct-digit-assignment.md) ·
[Solution capping](/rules/solution-capping.md) ·
[Solve cryptarithm](/usecases/solve-cryptarithm.md)

# Citations

- `internal/z07/cryptarithm.go` — `Solution`, `Result`, `parseEquation`, the no-zero set,
  `okSpace` shape checks, the permutation DFS, `spaceSize`, `capSol`.
