---
type: Domain
title: Cryptarithm puzzles (z07)
description: "Solve letter-for-digit equation puzzles over + - * / by pruned permutation search, honouring digit bijectivity, no leading zeros, and exact division, and report every solution."
tags: [domain, core, z07, combinatorics]
timestamp: "2026-10-27T16:32:44Z"
---

# Kind of subdomain

Core (the exercise itself). Cryptarithms: equations in which letters stand for digits, e.g.
`ABC * BD = EFGAH` solved as 623 * 29 = 18067.

# Rules

- Form: `<number> <op> <number> = <number>` with op ∈ `+ - * /`, letters A–Z only.
- Each letter maps to exactly one digit and no digit is used by two letters (a bijection over
  the letters present).
- A multi-digit number may not start with zero.
- Division is exact.
- All solutions must be findable — the exercise warns that some puzzles have several.

# Approach and pruning

Exhaustive search over digit permutations with pruning: leading letters skip digit 0, and a
preliminary length check rejects impossible shapes before any search (for `*` the result's
digit count must lie in [len(a)+len(b)−1, len(a)+len(b)], and for `/` the analogous relation
on the dividend). Multiplication overflow is guarded before comparison. The full search of
10P10 permutations costs a fraction of a second.

# Reporting all solutions

The full count of solutions is always reported, independent of the request. The requested
limit restricts only the returned list, and the response states that the list was truncated —
so "the puzzle has 476 solutions, here are 3 of them" is expressible. Within the library a
non-positive limit means *no limit*.

# Related

[Cryptarithm puzzle](/model/cryptarithm-puzzle.md) ·
[Distinct digit assignment](/invariants/distinct-digit-assignment.md) ·
[Solve cryptarithm](/usecases/solve-cryptarithm.md) ·
[Solution capping](/rules/solution-capping.md)

# Verified results

`ABC·BD=EFGAH` has 8 solutions (incl. 623·29=18067); `SEND+MORE=MONEY` has exactly one
(9567+1085=10652); `AB - CD = EF` has 476 solutions.

# Citations

- `internal/z07/cryptarithm.go` — equation parsing, letter set, no-zero letters, length
  pruning, permutation DFS with overflow guard, space size, full count with capped list.
- `internal/z07/cryptarithm_test.go`, `limit_test.go`, `leadingzero_test.go` — the verified
  puzzles, truncation semantics, leading-zero rule.
