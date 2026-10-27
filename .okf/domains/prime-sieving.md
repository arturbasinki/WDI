---
type: Domain
title: Prime sieving (z03)
description: "Sieve of Eratosthenes over three data structures (boolean slice, map, bitmap) with a cancellable, sequential timing comparison across increasing ranges."
tags: [domain, core, z03, number-theory, benchmark]
timestamp: "2026-10-27T16:32:44Z"
---

# Kind of subdomain

Core (the exercise itself). Sieve of Eratosthenes over [2, n], implemented over three
data structures, with the timings compared for growing n.

# The exercise is the comparison

The point is not primality — it is measuring how the choice of structure costs time at
10³, 10⁴, 10⁵, 10⁶, 10⁷, 10⁸. The three variants mirror the original Python structures:

1. a boolean slice ("list") — one flag per number,
2. a map of number → flag ("dictionary"),
3. a packed bitmap ("array") — eight flags per byte, the dense vectorised case.

The algorithm is identical in all three: mark candidates, cross out multiples of every
candidate p from p² up, the survivors are prime.

# Measurement honesty

Measurements are sequential (so they are comparable), with collection forced before each
measurement and collection left enabled so numbers reflect realistic costs. Cancellation is
checked *inside* a measurement, at a low-cost granularity, so cancelling a slow structure at
n = 10⁸ interrupts in a fraction of a second instead of after minutes
([Cancellation granularity](/rules/cancellation-granularity.md)). A measurement interrupted
mid-way is not reported as a row — partial rows would look like a genuinely slow structure.

All three structures must agree on the prime count for a given n; the map variant is
additionally ordered, because map iteration is not deterministic.

# Related

[Sieve benchmark row](/model/sieve-benchmark-row.md) ·
[Benchmark sieve](/usecases/benchmark-sieve.md) ·
[Measurement flow](/dataflows/measurement-flow.md) ·
[Cancellation inside measurement decision](/decisions/cancellation-inside-measurement.md)

# Verified results

Prime counts match reference values (168, 1229, 9592, 78498, 664579, 5761455). At n = 10⁷:
slice ≈ 46 ms, bitmap ≈ 45 ms, map ≈ 4 s.

# Citations

- `internal/z03/sieve.go` — the three sieve variants with cooperative cancellation, the
  benchmark loop, GC forcing, progress per completed measurement.
- `internal/z03/sieve_test.go`, `order_test.go`, `cancel_test.go` — reference counts, ordered
  map output, cancellation interrupting a single measurement.
