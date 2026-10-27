---
type: Actor
title: Examiner (task author)
description: The author of the exercise statements; their wording is the acceptance authority for each task domain.
tags: [actor, specification]
timestamp: "2026-10-27T16:32:44Z"
---

# Role

Not a user of the running system. The examiner wrote the exercise statements, and those
statements — including their example data and expected outputs — define what each task
domain must deliver. They are the party the product is ultimately answerable to.

# Binding expectations (as recorded in the repository)

- z01: solve N×N linear systems; must handle contradictory and underdetermined cases.
- z02: find as many Armstrong numbers as possible.
- z03: sieve with three data structures and compare timings for n = 10³…10⁸.
- z04: top-20 most frequent words, ignoring words shorter than 5 letters, for "Pan Tadeusz".
- z05: encrypt/decrypt with a numeric key over alphabet = space + 26 lowercase letters;
  break a given ciphertext with a dictionary; key length < 10.
- z06: plot y = f(x) over an interval, accounting for vertical asymptotes.
- z07: cryptarithms `<number> <op> <number> = <number>`, op ∈ + - * /; account for puzzles
  with multiple solutions.
- z08: sudoku puzzle with exactly one solution and as many empty cells as possible.

# Consequences for the system

Wording in the UI must not overclaim beyond what the implementation guarantees. The
"maximum number of empty cells" promise of z08 is deliberately softened in the task
catalog for exactly this reason — see [Task catalog invariants](/invariants/task-catalog-consistency.md).

# Citations

- `README.md` — the full exercise statements, kept verbatim under the solution description.
- `cmd/wdi-server/catalog_test.go` — `TestZ08TagBezObietnicyMaksimum` guards the over-claim.
