---
type: ValueObject
title: Cipher key
description: "A digit-string key whose digits supply the per-position shifts, cycling over its length; leading zeros are significant, length 1..9."
tags: [value-object, z05]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Value object. A finite sequence of decimal digits used as the shift schedule: the i-th
character of the message is shifted by the (i mod length)-th digit of the key.

# Identity

Compared as a **string of digits**, not as a number — leading zeros are significant
("007" is a valid, different key from "7"), and key enumeration in the attack includes them.

# Validity

- Digits `0–9` only; any other character is invalid.
- Length 1..9, from the exercise's "key length is less than 10".
- A digit of 0 is legal and means "no shift" for that position.

# Search space

Keys of length L number 10^L; enumerating lengths 1..9 is 111 111 111 keys, so the attack is
bounded by an explicit cap on keys tried rather than by the key syntax
([Cipher break search](/workflows/cipher-break-search.md)).

# Related

[Gronsfeld cipher](/domains/polyalphabetic-cipher.md) ·
[Cipher alphabet](/model/cipher-alphabet.md) ·
[Cipher text use case](/usecases/cipher-text.md)

# Citations

- `internal/z05/gronsfeld.go` — the key validation (digits, length 1..9) and the
  `key[i % len(key)]` indexing; the lexicographic key generator with zero padding.
- `internal/z05/gronsfeld_test.go` — reference vectors pinning the key indexing.
