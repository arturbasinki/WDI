---
type: UseCase
title: Solve a linear system
description: "Submit a linear system, get its classification, solution, RREF and the narrated elimination steps in one synchronous response."
tags: [use-case, z01, synchronous]
timestamp: "2026-10-27T16:32:44Z"
---

# Actor · goal

Learner — solve a system of N linear equations and see the reduction.

# Input

A text block: N on the first line, then N rows of N coefficients and the free term. The
interface also offers example data. Tokens may be separated by spaces, tabs, newlines or
commas; a sign written apart from its digits is accepted.

# Main flow

1. Submit the text.
2. Parse and validate synchronously — a malformed system is a 400 with the row and value
   named; no partial answer is produced ([Synchronous validation](/rules/synchronous-validation.md)).
3. Reduce with pivoted Gauss-Jordan, recording operations under the cell budget
   ([Animation step budget](/rules/animation-step-budget.md)).
4. Classify and answer in one response: status, values, RREF, rank, and the step narration.
5. The client animates the recorded steps; on truncation it replays what it received.

# Alternative flows

- **Contradictory system** → status `none`, one line stating there is no solution.
- **Underdetermined system** → status `infinite`, free unknowns named `t1…tk` and each pivot
  unknown expressed parametrically with the negated coefficient
  ([Solution status](/model/solution-status.md)).
- **Non-square / short / over-supplied input** → 400 stating how many numbers were expected
  and how many were found.
- **Non-finite coefficient** → 400, because a NaN answer would otherwise be presented as a
  unique solution.

# Verified outcome

The exercise's example yields x = 18, y = −58, z = 35. A 100×100 system answers in tens of
milliseconds at ~1.18 MB — asserted by the gate probe.

# Citations

- `internal/z01/solver.go` — `Parse`, `Solve`.
- `cmd/wdi-server/main.go` — the synchronous solve endpoint.
- `internal/z01/solver_test.go`, `param_test.go`, `stepsbudget_test.go`.
