---
type: ValueObject
title: Cipher alphabet
description: The 27-symbol cipher alphabet (space + 26 lowercase letters) whose ordering and size define the modular shift arithmetic.
tags: [value-object, z05]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Value object. The fixed 27-symbol symbol set the cipher operates on: a space followed by the
26 lowercase English letters, in that order.

# Meaning

Shift arithmetic is modular over this set: the position of a symbol in the alphabet is shifted
by a key digit and taken modulo 27. Its size (27) is the modulus; its ordering is what makes
"shift by k" meaningful, so the order is part of the identity of the cipher, not an
implementation detail.

# Consequences

- Anything outside the alphabet cannot be enciphered, which is why normalization drops it —
  and why that drop must be reported ([Normalization report](/model/normalization-report.md)).
- A plaintext that normalizes to nothing is an error, not an empty ciphertext.
- The same alphabet bounds what the dictionary attack can ever match.

# Related

[Gronsfeld cipher](/domains/polyalphabetic-cipher.md) ·
[Cipher key](/model/cipher-key.md) · [Cipher alphabet closure](/invariants/cipher-alphabet-closure.md)

# Citations

- `internal/z05/gronsfeld.go` — the alphabet constant and the modular shift arithmetic in
  encrypt, decrypt and the attack.
