---
type: Invariant
title: Cipher alphabet closure
description: Cipher input is confined to the 27-symbol alphabet; dropped characters are always reported and empty-after-normalization is an error.
tags: [invariant, z05, honesty]
timestamp: "2026-10-27T16:32:44Z"
---

# Statement

Everything the cipher processes is inside the 27-symbol alphabet (space + `a–z`). Text outside
it is never silently discarded: the set of dropped characters is reported, and a text that
normalizes to nothing is an error. Keys are digit strings of length 1..9.

# Rationale

The cipher is defined over that alphabet, so out-of-alphabet input has no encryption. The
failure mode observed in practice was silent: a Polish sentence returned 200 with an
almost-empty result, and the user had no way to learn that seven letters had been dropped. The
rule exists so that "what you see is what was enciphered".

# Enforcement

Normalization is a single, shared step at the entry of encrypt, decrypt and the attack; the
dropped set is sorted by code point and unique. The transport surfaces it only when non-empty,
and the client renders an explicit warning naming the characters.

# Related

[Normalization report](/model/normalization-report.md) ·
[Cipher alphabet](/model/cipher-alphabet.md) ·
[No silent truncation](/rules/no-silent-truncation.md)

# Citations

- `internal/z05/gronsfeld.go` — `normalizeText`, `validateKey` (digits, 1..9), `Alphabet`.
- `cmd/wdi-server/main.go` — `skipped` present only when non-empty.
- `cmd/wdi-server/web/app.js` — the dropped-characters warning in the z05 view.
