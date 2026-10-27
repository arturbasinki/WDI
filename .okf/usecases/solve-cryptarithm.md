---
type: UseCase
title: Solve a cryptarithm
description: "Submit a cryptarithm equation and receive every solution's count with a capped list of mappings, the full count always reported."
tags: [use-case, z07, synchronous]
timestamp: "2026-10-27T16:32:44Z"
---

# Actor · goal

Learner — solve a letter-for-digit puzzle and see every solution, not just the first.

# Main flow (synchronous)

1. Submit an equation `<number> <op> <number> = <number>` and an optional limit on returned
   solutions (default 50, ceiling 500, clamped above the ceiling; a negative limit is a 400
   explaining that 0 means "default" and −1 means "unlimited" at the library level).
2. Parse and check the puzzle's structure ([Cryptarithm puzzle](/model/cryptarithm-puzzle.md)).
3. Prune the search space by shape, then search digit permutations with the no-zero constraint
   ([Distinct digit assignment](/invariants/distinct-digit-assignment.md)).
4. Response: normalised formula, operator, letter count, space size, **full** solution count,
   `truncated`, up to `limit` solutions (each with its mapping and the digit-rendered
   equation), elapsed time.

# Alternative flows

- **Wrong token count, missing `=`, non-letter characters, an unknown operator** → 400 naming
  what was expected.
- **More than ten letters** → 400: ten digits cannot serve eleven letters.
- **Impossible shape** (e.g. a product whose digit count cannot match) → a successful answer
  with `count: 0`; "no solution" is data.

# Verified outcome

`ABC·BD=EFGAH` → 8 solutions incl. 623·29=18067; `SEND+MORE=MONEY` → exactly 9567+1085=10652;
`AB-CD=EF` → 476 solutions, of which a `limit=3` request returns 3 with `count: 476` and
`truncated: true`.

# Citations

- `internal/z07/cryptarithm.go` — `Solve`, `parseEquation`, pruning, DFS, `Count`/`Truncated`.
- `cmd/wdi-server/main.go` — the endpoint and its limit handling.
- `internal/z07/cryptarithm_test.go`, `limit_test.go`, `leadingzero_test.go`.
