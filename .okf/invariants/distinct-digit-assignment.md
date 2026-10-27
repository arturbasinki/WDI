---
type: Invariant
title: Distinct digit assignment
description: "Cryptarithm solutions are digit bijections: unique digits per letter, at most ten letters, no leading zero on multi-digit numbers, exact division."
tags: [invariant, z07]
timestamp: "2026-10-27T16:32:44Z"
---

# Statement

In a cryptarithm, the letter→digit assignment is a bijection over the letters present: each
letter takes exactly one digit, no digit is shared by two letters, at most ten distinct letters
may appear, and the first letter of any multi-digit number is never zero.

# Rationale

This is the definition of the puzzle. A repeated digit for two letters, or a leading zero,
produces an "answer" that is not a valid reading of the puzzle. The leading-zero rule in
particular must be applied during the search (as a pruning constraint), not by filtering
results afterwards, both for correctness of the count and for search cost.

# Corollaries

- `count` counts only assignments satisfying all of the above.
- Numbers are evaluated digit by digit, so a length-10 number is legal only if its first
  letter is not forced to zero by the same rule.
- Division is exact: `A / B = C` holds iff `B ≠ 0`, `A mod B = 0`, and `A / B = C`.

# Related

[Cryptarithm puzzle](/model/cryptarithm-puzzle.md) · [Cryptarithm](/domains/cryptarithm.md)

# Citations

- `internal/z07/cryptarithm.go` — the >10-letters rejection, the `noZero` set, the `used`
  digit array in the DFS, the exact-division check.
- `internal/z07/leadingzero_test.go` — leading-zero rule enforced.
- `internal/z07/cryptarithm_test.go` — verified solution counts.
