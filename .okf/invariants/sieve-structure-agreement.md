---
type: Invariant
title: Sieve structure agreement
description: "All three sieve structures yield the same ascending prime set for a given n, matching published reference counts."
tags: [invariant, z03, benchmark]
timestamp: "2026-10-27T16:32:44Z"
---

# Statement

For a given n, all three sieve structures must find the same set of primes — in particular the
same count, and in ascending order. The counts must match published reference values.

# Rationale

The exercise's deliverable is a timing comparison. If the structures disagreed about which
numbers are prime, the timings would compare three different computations. Ordering is part of
the invariant because the map variant iterates non-deterministically; without sorting, its
output order — and any test comparing lists — would be a coin flip.

# Reference values pinned by tests

n = 10³ → 168, 10⁴ → 1229, 10⁵ → 9592, 10⁶ → 78 498, 10⁷ → 664 579, 10⁸ → 5 761 455.

# Related

[Prime sieving](/domains/prime-sieving.md) · [Sieve benchmark row](/model/sieve-benchmark-row.md)

# Citations

- `internal/z03/sieve.go` — `SieveSliceCtx`, `SieveMapCtx` (with `sort.Ints`), `SieveBitCtx`.
- `internal/z03/sieve_test.go` — reference counts per structure.
- `internal/z03/order_test.go` — the map variant returns ascending output.
