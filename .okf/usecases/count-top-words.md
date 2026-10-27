---
type: UseCase
title: Count top words
description: "Request the top-N word ranking of the embedded canonical text or of a submitted text, with the applied N echoed back."
tags: [use-case, z04, synchronous]
timestamp: "2026-10-27T16:32:44Z"
---

# Actor · goal

Learner — see the most frequent words of "Pan Tadeusz", or of their own text.

# Main flow, canonical corpus

1. Request the ranking (no body). The default count of 20 is applied by the handler and echoed
   as `n` in the response ([Explicit defaults](/rules/explicit-defaults.md)).
2. The canonical text is taken from the content embedded in the artifact — no file on disk, no
   configuration ([Embedded corpus](/model/embedded-corpus.md)).
3. Response: the source marker, the ranked words, and the `n` used.

# Main flow, own text

1. Submit text and a count (1…100; outside that range the default 20 is applied and echoed).
2. Same analysis: case-folded, letters-only runs of ≥ 5
   ([Word length floor](/invariants/word-length-floor.md)).
3. Response: source marker `wlasny-tekst`, the ranking, the `n` used.

# Alternative flows

- A non-positive count reaching the library is an error, not a silent default.
- Ties: ranked alphabetically; the documented divergence from the original Python ranking is
  confined to tied positions.

# Verified outcome

Top 3: *rzekł* 155, *tylko* 149, *hrabia* 129.

# Citations

- `internal/z04/words.go` — `PanTadeusz`, `ExtractWords`, `TopWords`, `MinLength`, `DefaultTop`.
- `cmd/wdi-server/main.go` — the GET (canonical) and POST (own text) endpoints.
- `internal/z04/words_test.go` — explicit expected rankings.
