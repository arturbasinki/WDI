---
type: Decision
title: Uniqueness-proven digging for sudoku
description: "Generate sudoku by randomized fill then random-order digging, keeping a removal only when a count-to-two proves uniqueness; uniqueness is guaranteed, maximality is not."
tags: [decision, adr, z08, algorithm]
timestamp: "2026-10-27T16:32:44Z"
---

# Context

The exercise demands a sudoku with exactly one solution and the maximum number of empty cells.
Optimal digging is a hard combinatorial problem; the uniqueness requirement, however, is
decidable by exhaustive counting.

# Decision

Fill a complete grid by randomized backtracking, then **dig cells in random order, keeping a
removal only if counting the grid's solutions to a cap of two yields exactly one**. Repeat
until no cell can be removed. Prove uniqueness by counting, never by construction.

# Consequences

**Guaranteed:** uniqueness, for every seed, because every published grid passed an exhaustive
count with the same machinery that solves
([Unique solution](/invariants/sudoku-unique-solution.md),
[Sudoku solve workflow](/workflows/sudoku-solve.md)).

**Not guaranteed:** maximal blanks. The dig order is random and the procedure is monotonic, so
results are 55–58 blanks against a theoretical ceiling of ~64. The product states this in its
blurb and a test forbids the word "maximum"
([Task catalog consistency](/invariants/task-catalog-consistency.md)).

**Cost:** ~81 uniqueness counts per generation; with masked MRV counting that is 1–7 ms per
puzzle, so the cost is invisible to the learner.

# Rejected alternatives

Symmetry-aware digging and forced-chain-based clue minimisation (rejected as out of scope for
the exercise and for the standard-library-only constraint); fixed published puzzles (rejected:
the deliverable is a generator).

# Citations

- `internal/z08/sudoku.go` — `fillFull`, the dig loop with `CountSolutions(puzzle, 2)`,
  seeded generator, masked MRV counter.
- `internal/z08/sudoku_test.go`, `seed_test.go` — uniqueness and seed determinism.
- `docs/` — the external sudoku papers (multi-stage GA, SAT formulation, the 16-clue result)
  consulted for the uniqueness framing.
