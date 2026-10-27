---
type: Domain
title: Text word frequency (z04)
description: "Case-insensitive word-frequency ranking over a text with a five-letter floor, canonically the top 20 words of an embedded literary text."
tags: [domain, core, z04, text]
timestamp: "2026-10-27T16:32:44Z"
---

# Kind of subdomain

Core (the exercise itself). Report the most frequent words of a text, ignoring short words.

# Rules of the exercise

- Words shorter than five letters are not counted ("i", "lub", "się", "aby", "żeby" …).
- Counting is case-insensitive.
- The canonical answer is the top 20 words of Adam Mickiewicz's "Pan Tadeusz".
- Any user-supplied text may be analysed the same way.

# Decisions that make the answer stable

- A word is a maximal run of Unicode letters (an apostrophe does not split a word).
- Ranking: descending count, ties broken alphabetically. Ties are the documented divergence
  from the original Python implementation, whose `most_common` preserves insertion order;
  the non-tied ranking is guaranteed to agree.
- The requested number of results `k` must be positive; asking for zero is an error, not a
  silent substitution of the default. The default of 20 is applied explicitly by the caller
  ([Explicit defaults](/rules/explicit-defaults.md)).
- The canonical text is part of the artifact, so the default answer needs no file on disk
  ([Embedded content decision](/decisions/embedded-content.md)).

# Related

[Word tally](/model/word-tally.md) ·
[Word length floor](/invariants/word-length-floor.md) ·
[Count top words](/usecases/count-top-words.md) ·
[Embedded corpus](/model/embedded-corpus.md)

# Verified results

Top 3 for "Pan Tadeusz": *rzekł* 155, *tylko* 149, *hrabia* 129 — matching the Python
implementation.

# Citations

- `internal/z04/words.go` — embedded text, word extraction, counting, tie-stable ranking,
  the `k < 1` error, `MinLength`, `DefaultTop`.
- `internal/z04/pan-tadeusz.txt` — the canonical corpus.
- `internal/z04/words_test.go` — explicit expected rankings.
