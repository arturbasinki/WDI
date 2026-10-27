---
type: Workflow
title: Cipher break search workflow
description: "Enumerate digit keys by increasing length, accept the decryption whose words are over 80% dictionary-covered, with early abandon, key caps and cancellation."
tags: [workflow, z05, search]
timestamp: "2026-10-27T16:32:44Z"
---

# Scope

How the cipher-breaking search proceeds, and how it decides "found", "not found" and "stopped
early".

# Steps

1. Normalize the ciphertext; carry the dropped-character set through every outcome.
2. For key length L = 1 … maxLen, enumerate all 10^L digit strings lexicographically,
   **including leading zeros** ("007" is a distinct key).
3. For each candidate key: decrypt, then test the plaintext — split on whitespace, and count
   words absent from the dictionary.
   - Accept the key when more than 80% of the words are dictionary entries.
   - Abandon early as soon as more than 20% of the words examined so far have missed.
4. Count every key tried. Notify progress every 500 000 keys.
5. Stop on: a hit (return key + plaintext), the cap on keys tried (return "not found" with
   `capApplied` when the requested cap was itself clamped), exhaustion of all lengths,
   cancellation (return the cancellation with the partial count).

# Why the ratio rule works

The plaintext alphabet is only letters and spaces, so a correct key produces text whose words
are overwhelmingly dictionary words; a wrong key produces letter soup where the first few words
fail, and the early-abandon test cuts the cost per candidate dramatically.

# Known limits of the oracle

1 862 dictionary entries contain characters outside a–z and can never match, so a plaintext
using them cannot be recognised. The dictionary is content of the artifact
([Embedded corpus](/model/embedded-corpus.md)).

# Citations

- `internal/z05/gronsfeld.go` — `Break` (length loop, zero-padded key buffer, cap, notify),
  `breakKey` (the 80% accept / 20% abandon test).
- `internal/z05/gronsfeld_test.go` — the exercise proverb broken with key 142807.
