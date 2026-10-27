---
type: DataFlow
title: Corpus to ranking flow
description: "How embedded reference content reaches its consumers: the literary corpus through tokenising, filtering, tallying and ranking, and the dictionary as the cipher attack's oracle."
tags: [dataflow, z04, z05, content]
timestamp: "2026-10-27T16:32:44Z"
---

# Path

Reference text stored inside the artifact → token stream → filtered tally → ranked short list →
view; and, for the cipher, text → set-membership oracle → accept/abandon.

# Two consumers of embedded content

**Word ranking.** Embedded literary text → case-fold → maximal letter runs (apostrophe kept) →
drop runs shorter than five → count per distinct word → sort by count then alphabetically →
cut to the requested N → response with the N echoed
([Count top words](/usecases/count-top-words.md),
[Word length floor](/invariants/word-length-floor.md)).

**Cipher oracle.** Embedded dictionary → a set of headwords at load (duplicates collapse) →
for each candidate decryption, membership of each whitespace-delimited word → the
accept/abandon ratio decides the key ([Cipher break search](/workflows/cipher-break-search.md)).

# Properties of the source

- The content travels with the artifact: no path, no download, no configuration, so the
  canonical answer is reproducible on any machine
  ([Embedded content decision](/decisions/embedded-content.md)).
- A byte-order mark on the literary text is stripped on read.
- Dictionary entries containing characters outside `a–z` (1 862 of them) can never match a
  normalized plaintext, which bounds the attack's reachable answers.

# Related

[Embedded corpus](/model/embedded-corpus.md) ·
[Text frequency domain](/domains/text-frequency.md) ·
[Gronsfeld cipher domain](/domains/polyalphabetic-cipher.md)

# Citations

- `internal/z04/words.go` — embedded text, `ExtractWords`, `TopWords`.
- `internal/z05/gronsfeld.go` — the loaded dictionary set, `breakKey`.
- `internal/z04/pan-tadeusz.txt`, `internal/z05/slownik.txt`.
