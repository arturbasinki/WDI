---
type: ValueObject
title: Sieve benchmark row
description: "One measurement of one sieve data structure over one range: structure, n, prime count and duration — the row of the comparison table the exercise asks for."
tags: [value-object, z03, benchmark]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Value object. One cell of the comparison table the exercise demands: how one data structure
performed on one range.

# Schema

- `structure` — which variant: boolean slice ("list"), map ("dictionary"), bitmap ("array"),
- `n` — the range end,
- `count` — how many primes were found,
- `seconds` — the measured wall-clock duration.

Identity for comparison purposes is the pair (structure, n); the pair must appear at most once
per benchmark run.

# Semantics

- The count is a correctness signal, not just a by-product: all three structures must agree
  for the same n, and that agreement is what makes the timing comparison meaningful.
- A measurement interrupted by cancellation never yields a row. Reporting a truncated
  measurement as if it were a real one would make a structure look slow for the wrong reason.
- Measurements are taken sequentially, with collection forced beforehand, so rows are
  comparable with each other.

# Presentation

The client draws the rows as bars on a **logarithmic** scale (fastest structure 0%, slowest
100%), because a linear scale flattens everything against the map at 10⁸; sub-millisecond
durations are rendered as "< 1 ms".

# Related

[Prime sieving](/domains/prime-sieving.md) · [Benchmark sieve](/usecases/benchmark-sieve.md) ·
[Measurement flow](/dataflows/measurement-flow.md)

# Citations

- `internal/z03/sieve.go` — `Row`, the sequential measurement loop, collection forcing,
  the no-partial-row-on-cancel rule.
- `cmd/wdi-server/web/app.js` — `fmtSec`, the logarithmic bar helpers in the z03 view.
