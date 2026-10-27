---
type: UseCase
title: Break the cipher
description: "Run a cancellable dictionary attack against a Gronsfeld ciphertext, poll keys-tried progress, and receive the key, plaintext, keys tried and cap status."
tags: [use-case, z05, async]
timestamp: "2026-10-27T16:32:44Z"
---

# Actor · goal

Learner — break a Gronsfeld ciphertext with a dictionary attack and see the key being tried out.

# Main flow (background job)

1. Submit the ciphertext, an optional maximum key length (default 6), an optional cap on keys
   tried (default 2 000 000).
2. Synchronous validation: ciphertext non-empty, `maxLen` 1…9, `capKeys` 1…100 000 000;
   otherwise 400 and no job.
3. Take a computation slot (else 429), receive a job id.
4. Poll: progress carries keys tried, the key length currently exhausted, `done` (never true
   mid-attack).
5. Final result: found flag, the key, the recovered plaintext, keys tried, elapsed
   milliseconds, any `skipped` characters, and `capApplied`
   ([Cipher break search workflow](/workflows/cipher-break-search.md)).

# Alternative flows

- **Cap clamped** (a request above 10⁸) → 400 at the edge; a cap accepted and then reached
  mid-attack → a completed job with `found: false` and `capApplied: true`, which is how the
  client distinguishes "stopped at the cap" from "searched everything".
- **Cancel** → the job ends `cancelled`; the loop checks cancellation per candidate key.
- **No key found within the cap** → a normal, successful answer of "not found" with the keys
  tried — a negative result is data, not an error.

# Verified outcome

The exercise's proverb is broken with key **142807** → "there is no use crying over spilt
milk", after 253 918 keys tried.

# Citations

- `internal/z05/gronsfeld.go` — `Break`, `breakKey`, the key generator, the 80%/20% ratio rule,
  the cap and `CapApplied`.
- `cmd/wdi-server/main.go` — the break endpoint and its validation.
- `internal/z05/gronsfeld_test.go` — the attack on the exercise ciphertext.
