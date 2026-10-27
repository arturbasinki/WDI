---
type: BusinessRule
title: Solution capping
description: A solution limit trims the returned list only; the full solution count and the truncation flag are always reported.
tags: [business-rule, z07, honesty]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

A limit on returned solutions restricts **the list only**. The full count of solutions is always
reported, and truncation is stated. Asking for fewer solutions than exist must never look like
"the puzzle has that many solutions".

# Numbers in force

- Web layer: default 50 returned solutions, ceiling 500 — above the ceiling the value is
  clamped, not rejected, and a negative value is rejected with an explanation.
- Library: `limit <= 0` means no limit; the count is independent of the cap in every case.

# Worked example

`AB - CD = EF` has 476 solutions. A request with `limit = 3` returns 3 entries with
`count: 476` and `truncated: true`.

# Why

The exercise explicitly warns that puzzles may have more than one solution. A truncated list
presented as complete is the exact defect the exercise names. Counting is exhaustive anyway, so
the full count costs nothing extra to report.

# Related

[Cryptarithm puzzle](/model/cryptarithm-puzzle.md) ·
[Layer ownership of bounds](/rules/layer-ownership-of-bounds.md) ·
[No silent truncation](/rules/no-silent-truncation.md) ·
[Solve cryptarithm](/usecases/solve-cryptarithm.md)

# Citations

- `internal/z07/cryptarithm.go` — `capSol` (with `math.MaxInt` for no limit), `Count`,
  `Truncated`.
- `cmd/wdi-server/main.go` — default 50, ceiling 500, rejection of a negative limit.
- `internal/z07/limit_test.go` — full count with a capped list.
