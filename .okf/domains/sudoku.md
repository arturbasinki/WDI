---
type: Domain
title: Sudoku generation and solving (z08)
description: "Generate sudoku puzzles with a provably unique solution and many empty cells, by random fill plus uniqueness-preserving dig, and solve/count any grid."
tags: [domain, core, z08, combinatorics]
timestamp: "2026-10-27T16:32:44Z"
---

# Kind of subdomain

Core (the exercise itself). Generate sudoku puzzles that have exactly one solution and as
many empty cells as possible; also solve any supplied grid.

# Generation, in two phases

1. **Fill.** Complete a full 9×9 grid by backtracking with candidates in random order — the
   base is guaranteed to be a valid complete game.
2. **Dig.** Visit cells in random order; remove a value and count the grid's solutions with a
   counter that stops at two. If uniqueness survives, the removal stays; otherwise the value
   is restored. Repeat until no cell can be removed.

The result is a puzzle with a proven unique solution and many empty cells — typically 55–58
(verified for seeds 1–8 and 42). The theoretical ceiling for a unique-solution puzzle is
higher (~64), so the promise is stated as "many", never as "maximum": digging is monotonic,
not optimal.

# Counting and solving

Backtracking with the MRV heuristic (next cell = the one with the fewest candidates) and
bit masks for row, column and box. A typical hard puzzle resolves in tens of milliseconds.
Counting is capped, and a cap hit is reported as such; a grid that is internally contradictory
counts as zero solutions.

# Determinism

Generation is driven by an explicit seed, so a puzzle is reproducible; a seed of zero means
"take the current time".

# Related

[Sudoku puzzle](/model/sudoku-puzzle.md) ·
[Unique solution invariant](/invariants/sudoku-unique-solution.md) ·
[Generate sudoku](/usecases/generate-sudoku.md) · [Solve sudoku](/usecases/solve-sudoku.md) ·
[Sudoku solve workflow](/workflows/sudoku-solve.md) ·
[Unique-solution digging decision](/decisions/unique-solution-digging.md)

# Citations

- `internal/z08/sudoku.go` — grid type, validation, masked MRV counting with cap, random
  full fill, dig-while-unique loop, first-solution recovery.
- `internal/z08/sudoku_test.go`, `seed_test.go` — uniqueness verified by the counter,
  seed determinism.
- `docs/` — three external sudoku research papers consulted for the uniqueness/SAT framing.
