---
type: Aggregate
title: Sudoku puzzle
description: "Aggregate of a 81-cell sudoku position: generated puzzle with its proven-unique solution and blank count, plus a solve/count result for arbitrary grids."
tags: [aggregate, z08]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Aggregate. A sudoku position: 81 cells (0 = empty) with the row/column/box constraints, in two
roles — a generated puzzle (with its proven-unique solution) and an arbitrary grid supplied by
a learner for solving.

# Schema

Generated puzzle: `puzzle` (the grid with blanks), `solution` (the complete grid it came from),
`blanks` (empty cells), `removed` (cells successfully dug), `timeMs`.
Solve result: `solved`, `solution` (the first found), `count` (solutions counted up to the cap),
`capped` (counting stopped at the cap), `timeMs`.

# Invariants

- A grid is 81 values, each 0..9 ([Grid validation](/invariants/sudoku-grid-validation.md)).
- A generated puzzle has **exactly one** solution — enforced by counting to two after every
  dig and restoring the value when uniqueness is lost
  ([Unique solution](/invariants/sudoku-unique-solution.md)).
- A grid that already contradicts itself (a repeated digit in a row, column or box) counts as
  zero solutions.
- Generation is deterministic for a given seed.

# Candidate mechanics

Rows, columns and boxes are tracked as bit masks over digits 1..9; candidates for a cell are
the complement of the union of its three masks. The next cell is chosen by MRV (fewest
candidates, stopping early at a cell with exactly one), which is what makes a hard puzzle
resolve in tens of milliseconds.

# Related

[Sudoku](/domains/sudoku.md) · [Generate sudoku](/usecases/generate-sudoku.md) ·
[Solve sudoku](/usecases/solve-sudoku.md) · [Sudoku solve workflow](/workflows/sudoku-solve.md)

# Citations

- `internal/z08/sudoku.go` — `Grid`, `Result`, `SolveResult`, `ValidateGrid`,
  `CountSolutions` (masks, MRV, cap), `fillFull`, the dig loop, `fillFirst`.
- `cmd/wdi-server/web/app.js` — the single z08 board: generate, enter own puzzle, check, solve.
