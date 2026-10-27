---
type: Domain
title: Gronsfeld cipher (z05)
description: "Gronsfeld polyalphabetic cipher with a numeric key over a 27-symbol alphabet, its inverse, and a dictionary-attack breaker over enumerated keys."
tags: [domain, core, z05, cryptography]
timestamp: "2026-10-27T16:32:44Z"
---

# Kind of subdomain

Core (the exercise itself). The Gronsfeld cipher: a Caesar variant whose shift sequence comes
from the digits of a numeric key, plus breaking a ciphertext by dictionary attack.

# Cipher rules

- The alphabet is exactly 27 symbols: a space followed by the 26 lowercase English letters.
- Successive characters of the message are shifted by successive digits of the key, cycling
  over the key length. Key 312063 over "computer science" yields "fpop whsbsilhoee".
- The key is digits only, length 1..9 (the exercise states key length below 10).

# Normalization must be visible

Input is lower-cased and everything outside the alphabet is dropped. That must not be silent:
the result carries the set of unique dropped characters, sorted by code point, so that a
Polish sentence reports exactly which letters were removed from the answer. If nothing
remains after normalization, that is an error.

# Dictionary attack

Candidate keys are enumerated by increasing length, lexicographically within a length
(including leading zeros, e.g. "007"). Each candidate decrypts the ciphertext and the
plaintext is accepted when more than 80% of its words are dictionary entries — with an early
abandon as soon as more than 20% of the words examined have failed. The dictionary is
embedded (25 991 lines, 25 839 unique headwords); entries containing characters outside
a–z can never match, because the plaintext alphabet excludes them.

Total keys tried is bounded by a caller-supplied cap; a cap above the ceiling is clamped and
the result states that the cap was applied, so "stopped at the cap" is distinguishable from
"exhausted the key space".

# Related

[Cipher key](/model/cipher-key.md) · [Cipher alphabet](/model/cipher-alphabet.md) ·
[Normalization report](/model/normalization-report.md) ·
[Cipher break search workflow](/workflows/cipher-break-search.md) ·
[Encrypt/decrypt text](/usecases/cipher-text.md) ·
[Break cipher](/usecases/break-cipher.md)

# Verified results

The exercise's proverb is broken with key **142807** → "there is no use crying over spilt
milk" (253 918 keys tried).

# Citations

- `internal/z05/gronsfeld.go` — alphabet, key validation, normalization with the dropped set,
  encrypt/decrypt, key generator, dictionary match ratio, capped break with progress.
- `internal/z05/slownik.txt` — the embedded dictionary.
- `internal/z05/gronsfeld_test.go` — exact reference vectors in both directions; a key-index
  off-by-one is caught by them.
