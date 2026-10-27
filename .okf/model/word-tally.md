---
type: ValueObject
title: Word tally
description: "A word and its occurrence count in an analysed text, the unit of the frequency ranking with case folding and a length floor."
tags: [value-object, z04]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Value object: a word paired with the number of its occurrences. Compared by value, no identity.

# Schema

- `word` — the normalised word as counted,
- `count` — its occurrences in the analysed text.

# Rules of formation

- Lower-cased; a maximal run of Unicode letters (an apostrophe does not break a word).
- Words shorter than five letters never enter the tally
  ([Word length floor](/invariants/word-length-floor.md)).
- Ranking: descending `count`, ties alphabetically by `word`, then cut to the requested number.
- The tally is over a *set* of distinct words, so the ranking of tied words is deterministic;
  the documented divergence from the original Python ranking is confined to tied positions.

# Related

[Text word frequency](/domains/text-frequency.md) · [Word length floor](/invariants/word-length-floor.md)

# Citations

- `internal/z04/words.go` — `WordCount`, `ExtractWords`, `TopWords` and its comparator.
- `internal/z04/words_test.go` — explicit expected lists with tie ordering.
