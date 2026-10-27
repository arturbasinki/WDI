---
type: ValueObject
title: Normalization report
description: "The outcome of preparing text for the cipher: the encipherable text plus the sorted set of unique characters dropped, so nothing disappears silently."
tags: [value-object, z05, honesty]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Value object. What normalization did to the supplied text: the encipherable text plus the
**set of unique characters that were removed**.

# Normalization

Lower-case, then drop every character outside the cipher alphabet (i.e. keep only `a–z` and
space). Dropped characters are collected **uniquely** and sorted by code point, reported as
they appear in the normalized (lower-cased) form — that is the alphabet the cipher works on.

# Why it is a first-class object

Silent normalization destroyed trust: "ZAŻÓŁĆ GĘŚLĄ" returned HTTP 200 with a near-empty
result and seven Polish letters vanished without a trace. Now the report is part of every
cipher answer:
- encryption/decryption return the skipped set (omitted entirely when nothing was dropped),
- the interface shows a warning naming the dropped characters,
- the attack carries the same set on its result,
- if nothing remains after normalization, that is an error rather than an empty success.

# Related

[Cipher alphabet](/model/cipher-alphabet.md) ·
[No silent truncation](/rules/no-silent-truncation.md) ·
[Cipher alphabet closure](/invariants/cipher-alphabet-closure.md)

# Citations

- `internal/z05/gronsfeld.go` — `normalizeText` returning (text, skipped, error) with the
  code-point sort.
- `cmd/wdi-server/main.go` — the cipher handler includes `skipped` only when non-empty.
- `cmd/wdi-server/web/app.js` — the z05 view renders the "pominięte znaki" warning.
