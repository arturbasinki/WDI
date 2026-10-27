---
type: UseCase
title: Check / solve a sudoku grid
description: "Submit a sudoku grid to count its solutions up to a cap and, on request, receive the first solution for display."
tags: [use-case, z08, synchronous]
timestamp: "2026-10-27T16:32:44Z"
---

# Actor · goal

Learner — check the puzzle they typed (without being shown the answer) or have it solved.

# Main flow (synchronous)

1. Submit an 81-value grid (0 = empty) — validated for shape
   ([Sudoku grid validation](/invariants/sudoku-grid-validation.md)).
2. Count the grid's solutions with the masked, MRV-guided search, stopping at a cap (100 for
   this use case).
3. **Check**: the client presents solvability and the count without revealing the completed
   grid.
4. **Solve**: the client fills the whole board from the returned first solution.
5. Response: `solved`, `solution` (present only when solvable), `count`, `capped` (counting hit
   the cap), elapsed time.

# Alternative flows

- **Malformed grid** → 400 naming the cell.
- **Self-contradictory grid** → `solved: false`, `count: 0`; not a transport error.
- **Many solutions** → `count` up to the cap with `capped: true`, so "at least 100" is
  expressible rather than a misleading exact number.

# Related

The counting machinery is the same one the generator uses to prove uniqueness, with a cap of 2.

# Citations

- `internal/z08/sudoku.go` — `Solve`, `CountSolutions`, `fillFirst`, the cap and `Capped`.
- `cmd/wdi-server/main.go` — the solve endpoint (cap 100).
- `cmd/wdi-server/z08_ui_test.go` — the single-board check/solve behaviour.
