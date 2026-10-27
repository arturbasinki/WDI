---
type: UseCase
title: Encipher / decipher text
description: "Encipher or decipher text with a digit key over the 27-symbol alphabet, reporting any characters normalization dropped."
tags: [use-case, z05, synchronous]
timestamp: "2026-10-27T16:32:44Z"
---

# Actor · goal

Learner — encipher or decipher a text with a numeric key, and be told what the cipher could
not represent.

# Main flow

1. Submit mode (`enc` or `dec`), text, key.
2. Validate the mode synchronously; an unknown mode is a 400.
3. Normalize the text (lower-case, drop out-of-alphabet symbols) collecting the dropped set,
   and validate the key (digits, length 1..9).
4. Apply the modular shift over the 27-symbol alphabet, cycling the key.
5. Response: the result, plus `skipped` **only when** characters were dropped.
6. The client shows a warning naming the dropped characters.

# Alternative flows

- **Nothing encipherable** (e.g. a text of only digits and punctuation) → 400 stating that the
  text contains no letters a–z.
- **Key of the wrong length or with non-digits** → 400.
- **Polish input** → answer plus the exact list of removed letters, e.g. "ZAŻÓŁĆ GĘŚLĄ" →
  result `" bahm"` with skipped `ó ą ć ę ł ś ż`
  ([Cipher alphabet closure](/invariants/cipher-alphabet-closure.md)).

# Verified outcome

Key 312063 over "computer science" → "fpop whsbsilhoee", and the inverse.

# Citations

- `internal/z05/gronsfeld.go` — `Encrypt`, `Decrypt`, `normalizeText`, `validateKey`.
- `cmd/wdi-server/main.go` — the cipher endpoint and the conditional `skipped` field.
- `internal/z05/gronsfeld_test.go` — exact vectors both directions.
