---
type: Workflow
title: Sudoku solve / count workflow
description: "Mask-guided, MRV-ordered backtracking used both to count solutions up to a cap and to find the first solution of a sudoku grid."
tags: [workflow, z08, search]
timestamp: "2026-10-27T16:32:44Z"
---

# Scope

The shared backtracking machinery behind sudoku generation, counting and solving.

# Steps

1. **Seed masks.** From the given grid, build a bit mask per row, per column and per box over
   digits 1..9. If a digit is already repeated in any group, the grid is contradictory → the
   count is zero, immediately.
2. **Select the next cell (MRV).** Among empty cells, take the one with the fewest candidates,
   where candidates = the complement of (row ∪ column ∪ box masks). Stop scanning early at a
   cell with exactly one candidate. A cell with zero candidates fails the branch.
3. **Try each candidate**, applying and undoing the three masks, and recurse.
4. **Terminate** when no empty cell remains: that is one solution. Counting stops as soon as
   the cap is reached; solving returns the first complete grid.

# Two consumers of the same machinery

- **Counting, cap 2** — used by the generator after every dig: exactly one solution keeps the
  removal, anything else restores the value
  ([Unique solution](/invariants/sudoku-unique-solution.md)).
- **Counting, cap 100** — used by the solve use case, with the cap hit reported as `capped`
  ([Check / solve a sudoku grid](/usecases/solve-sudoku.md)).

# Why it is fast enough

Masks make candidate computation O(1) per cell, and MRV picks the most constrained cell, which
prunes the tree aggressively — a typical hard puzzle resolves in tens of milliseconds.
Generation (one full fill plus ~81 uniqueness counts) costs 1–7 ms.

# Related

[Sudoku puzzle](/model/sudoku-puzzle.md) · [Sudoku](/domains/sudoku.md) ·
[Unique-solution digging decision](/decisions/unique-solution-digging.md)

# Citations

- `internal/z08/sudoku.go` — `CountSolutions` (masks, MRV, cap), `fillFull` (randomized fill),
  `fillFirst`, `popcount16`, `boxOf`.
