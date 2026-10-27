---
type: Invariant
title: Sudoku grid validation
description: "A sudoku grid is exactly 81 values each in 0..9, and solution counting requires a cap of at least 1."
tags: [invariant, z08]
timestamp: "2026-10-27T16:32:44Z"
---

# Statement

A grid submitted for solving must be exactly 81 values, each in 0..9. A cap on solution
counting must be at least 1.

# Rationale

81 is the board, 0 encodes "empty", 1..9 are the digits. Anything else is a client bug and must
be rejected with a message naming the offending cell (1-based), rather than being solved into a
nonsense answer.

# Related rule (not a shape check)

A grid that is well-shaped but self-contradictory — the same digit twice in a row, column or
box — is not rejected: it is a legitimate position with zero solutions, and the count reports
that.

# Related

[Sudoku puzzle](/model/sudoku-puzzle.md) · [Sudoku](/domains/sudoku.md) ·
[Error response contract](/contracts/error-response.md)

# Citations

- `internal/z08/sudoku.go` — `ValidateGrid` (length and per-cell range), the `cap < 1` check,
  the contradiction short-circuit to zero in `CountSolutions`.
- `internal/z08/sudoku_test.go` — validation cases.
