---
type: Invariant
title: Word length floor
description: Only letter runs of at least five characters (case-folded) enter the word-frequency tally.
tags: [invariant, z04]
timestamp: "2026-10-27T16:32:44Z"
---

# Statement

A word enters the frequency tally only if it is a maximal run of letters of length ≥ 5 after
case folding. Words shorter than five letters are never counted, and a run is broken by any
non-letter except an apostrophe.

# Rationale

Stated by the exercise: short Polish function words ("i", "lub", "się", "aby", "żeby") would
dominate any ranking and drown the signal. Five letters is the exercise's own threshold, so it
is a domain rule, not a tuning knob.

# Corollaries

- The floor is applied to the **word**, not to the text.
- Case folding happens before counting, so "Hrabia" and "hrabia" are one word.
- Because an apostrophe is kept inside a word, French-style text does not fragment; in Polish
  the choice is immaterial, which is exactly why it is documented.
- The threshold is a constant of the domain (5), exposed by name so a re-implementation cannot
  drift silently.

# Related

[Text word frequency](/domains/text-frequency.md) · [Word tally](/model/word-tally.md)

# Citations

- `internal/z04/words.go` — `MinLength = 5`, `ExtractWords` (letter runs, apostrophe, flush).
- `internal/z04/words_test.go` — short words excluded, counts of the canonical text.
