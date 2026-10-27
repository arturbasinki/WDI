---
type: Entity
title: Embedded corpus
description: Reference text shipped inside the deployed artifact — the canonical literary corpus and the English dictionary — identified by name and immutable at run time.
tags: [entity, z04, z05, content]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Entity. A body of reference text that ships **inside** the artifact and is identified by its
name, not by its content bytes on disk.

# Instances

| Name | Content | Serves |
|---|---|---|
| Canonical literary text | "Pan Tadeusz" (BOM stripped on read) | the default corpus for [Text word frequency](/domains/text-frequency.md) |
| English dictionary | 25 991 lines = 25 839 unique headwords | the breaking oracle of [Gronsfeld cipher](/domains/polyalphabetic-cipher.md) |

# Rules

- The corpus is content of the module that analyses it; there is no file path, no download,
  no configuration pointing at it ([Embedded content decision](/decisions/embedded-content.md)).
- The dictionary is a set, not a list: duplicates collapse at load.
- Dictionary entries containing characters outside `a–z` can never match a decrypted
  plaintext, because the cipher alphabet excludes them (1 862 such entries). This is a
  property of the corpus, and it bounds what the attack can find.
- Because the corpus is embedded, the default answer to the word-frequency exercise is
  reproducible on any machine with no preparation.

# Related

[Text frequency flow](/dataflows/corpus-to-ranking-flow.md) ·
[Embedded content decision](/decisions/embedded-content.md)

# Citations

- `internal/z04/words.go` — the embedded text, `PanTadeusz()` stripping the BOM.
- `internal/z04/pan-tadeusz.txt` — the literary corpus.
- `internal/z05/gronsfeld.go` — the embedded dictionary loaded as a set.
- `internal/z05/slownik.txt` — the dictionary corpus.
