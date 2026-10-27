---
type: ValueObject
title: Search progress
description: "The pollable state of a long computation: processed/limit, results so far as a sorted copy, and a done flag reserved for genuinely finished runs."
tags: [value-object, progress, jobs]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Value object. The observable state of a long search at one moment in time.

# Schema (narcissistic-number search)

- `processed` — how many candidates have been examined,
- `limit` — the upper bound of the searched range,
- `found` — numbers found so far, **sorted ascending**, always a copy (never aliasing the
  searching state),
- `done` — true only for the final state of a completed search.

Analogous progress shapes exist for the sieve benchmark (completed measurements / total, plus
the last row) and for the cipher break (keys tried, current key length, found flag, plaintext,
key, done).

# Rules carried by the shape

- Progress is emitted at a steady cadence, not per candidate: ≈10 Hz for the number search,
  every 500 000 keys tried for the cipher attack, once per completed measurement for the
  benchmark ([Progress cadence](/rules/progress-cadence.md)).
- `done: true` is reserved for a genuinely finished run — never emitted mid-attack.
- A cancelled run returns the partial progress **with** the cancellation; it is not presented
  as a completed result.
- The progress producer must be joined before the final result is assembled, otherwise the
  final ordering races with the progress reader.
- A nil progress callback is accepted everywhere (treated as "no progress").

# Related

[Job](/model/job.md) · [Narcissistic numbers](/domains/narcissistic-numbers.md) ·
[Prime sieving](/domains/prime-sieving.md) · [Gronsfeld cipher](/domains/polyalphabetic-cipher.md)

# Citations

- `internal/z02/armstrong.go` — `Progress`, ticker goroutine, join-before-final, sorted copies.
- `internal/z03/sieve.go` — `Progress{Done,Total,Last}` per measurement.
- `internal/z05/gronsfeld.go` — `BreakProgress`, the 500 000-key notification cadence.
