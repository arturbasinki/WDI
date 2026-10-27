---
type: Decision
title: Cancellation inside the measured loop
description: "Check cancellation inside sieve loops once per 2^20 iterations (mask test), including the initialisation pass, so the slowest measurement is interruptible in a fraction of a second."
tags: [decision, adr, cancellation, z03]
timestamp: "2026-10-27T16:32:44Z"
---

# Context

The sieve benchmark's slowest cell is the map structure at n = 10⁸ — 4…15 s per measurement,
~7 s under the race detector for the initialisation pass alone. Cancellation was checked only
between measurements, so cancelling a benchmark took seconds to minutes, and a single
measurement could not be interrupted at all.

# Decision

Check cancellation **inside** the loops of each sieve, once every 2²⁰ iterations, using a
full-bit mask so the test is one comparison amortised over a million iterations. Put the check
in every loop, including the table-initialisation pass, because at n = 10⁸ initialisation alone
is long enough to be the dominant cost.

# Consequences

Measured: cancelling went from ~7 s to ~0.04 s, and the cost of the check is unmeasurable in
the timings. The rule "a cancelled measurement produces no row" is required alongside, so
partial measurements never enter the comparison table
([Sieve benchmark row](/model/sieve-benchmark-row.md)).

# Generality

The rule is not specific to the sieve: any computation behind the
[cancellable port](/contracts/solver-package-port.md) must notice cancellation at a granularity
small enough that the slowest unit of work is interruptible
([Cancellation granularity](/rules/cancellation-granularity.md)).

# Related

[Prime sieving](/domains/prime-sieving.md) · [Benchmark sieve](/usecases/benchmark-sieve.md) ·
[Computation cancelled](/events/computation-cancelled.md)

# Citations

- `internal/z03/sieve.go` — `ctxTick = 1 << 20`, the `canceled` helper, checks in
  initialisation, crossing-out and collection loops.
- `internal/z03/cancel_test.go` — a single measurement is interrupted.
- `RAPORT-JAKOSC-KODU.md` §5 P2 item 14 — the original finding.
