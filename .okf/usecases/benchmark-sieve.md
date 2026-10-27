---
type: UseCase
title: Benchmark the sieve
description: "Run a cancellable benchmark of three sieve structures over chosen ranges, poll per-measurement progress, and receive the comparison table."
tags: [use-case, z03, async, benchmark]
timestamp: "2026-10-27T16:32:44Z"
---

# Actor · goal

Learner — compare how three data structures perform at sieving primes, for a chosen set of
ranges, with visible progress.

# Main flow (background job)

1. Submit a list of `n` values; an empty list means the exercise's ladder 10³, 10⁴, 10⁵, 10⁶,
   10⁷, 10⁸.
2. Synchronous validation: every n in 1000 … 100 000 000; one bad value rejects the whole
   request with 400 and no job ([Synchronous validation](/rules/synchronous-validation.md)).
3. Take a computation slot (else 429), receive 202 with a job id.
4. Poll: progress reports completed measurements out of `len(ns) × 3`, plus the last row
   measured.
5. Final result: the rows of the comparison table
   ([Sieve benchmark row](/model/sieve-benchmark-row.md)).
6. The client draws the rows as bars on a logarithmic scale, so the map at 10⁸ does not flatten
   everything else.

# Alternative flows

- **Cancel** → the running measurement is interrupted inside, in a fraction of a second
  ([Cancellation granularity](/rules/cancellation-granularity.md)); the job ends `cancelled`
  with the rows completed so far. No partial measurement is reported as a row.
- **Out-of-range n** → 400 before any job exists.

# Acceptance

Counts agree across the three structures and match reference values
([Sieve structure agreement](/invariants/sieve-structure-agreement.md)).

# Citations

- `internal/z03/sieve.go` — `Benchmark`, `Progress`, the three structures.
- `cmd/wdi-server/main.go` — the benchmark endpoint.
- `internal/z03/cancel_test.go`, `sieve_test.go`, `order_test.go`;
  `cmd/wdi-server/handlers_test.go`.
