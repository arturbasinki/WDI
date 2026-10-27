---
type: Invariant
title: Sudoku unique solution
description: "Every generated sudoku has exactly one solution, verified by counting to two after each cell removal."
tags: [invariant, z08]
timestamp: "2026-10-27T16:32:44Z"
---

# Statement

Every generated sudoku puzzle has **exactly one** solution. Generation never publishes a grid
whose solution count differs from one.

# How it is guaranteed

Not by construction, but by verification at every step: each dig is followed by counting the
grid's solutions with a counter that stops at two, and the removed value is restored if the
count is not exactly one. Generation ends when no cell can be removed without breaking
uniqueness.

# Consequences

- The guarantee is a proof by exhaustive count, so it holds for every seed — including the
  ones presented in the documentation.
- "As many empty cells as possible" is *not* guaranteed: the dig order is random and the
  procedure is monotonic, so the result is 55–58 blanks, while the theoretical ceiling for a
  unique puzzle is higher (~64). The product's blurb states this honestly, and a test forbids
  the word "maximum" in the sudoku blurb.
- Uniqueness is checked with the same machinery that solves, so a bug in the counter would
  corrupt generation; the counter is verified against independently checked puzzles.

# Related

[Sudoku](/domains/sudoku.md) · [Sudoku puzzle](/model/sudoku-puzzle.md) ·
[Unique-solution digging decision](/decisions/unique-solution-digging.md) ·
[Task catalog invariants](/invariants/task-catalog-consistency.md)

# Citations

- `internal/z08/sudoku.go` — the dig loop: `CountSolutions(puzzle, 2) != 1` ⇒ restore.
- `internal/z08/sudoku_test.go`, `seed_test.go` — uniqueness by the counter, seeded grids.
- `cmd/wdi-server/catalog_test.go` — the blurb may not promise a maximum.
