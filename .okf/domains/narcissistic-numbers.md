---
type: Domain
title: Narcissistic numbers (z02)
description: "Exhaustive, parallel, incrementally optimised search for numbers equal to the sum of their digits raised to the digit-count power, to a verified reach of 10^12."
tags: [domain, core, z02, number-theory]
timestamp: "2026-10-27T16:32:44Z"
---

# Kind of subdomain

Core (the exercise itself). Find as many narcissistic/Armstrong numbers as possible: an
N-digit natural number equal to the sum of its digits each raised to the power N
(153 = 1³+5³+3³).

# Approach and its intent

Exhaustive search over [0, limit] with an incremental digit-power sum, parallelised by
splitting the range into fixed-size blocks consumed by a pool of workers sized to the
available cores.

The central optimisation works on blocks of ten consecutive numbers sharing the same leading
part: the digit-power sum of the leading part is computed once per block, and the ten
candidates differ only by `pow[digit]`. One division, one addition and one comparison per
number — roughly an order of magnitude faster than recomputing the digit sum per number.
Ranges are grouped per digit count so the power is constant inside a block.

# Scale as the point of the exercise

The exercise asks for "as many as possible", so reach is the deliverable: the verified sweep
is to 10¹² (31 minutes on 12 cores) producing 41 numbers, all listed and pinned by tests;
values up to 10⁹ are additionally covered by unit tests.

# Honest boundaries

The search accepts any non-negative limit — searching [0,100] is correct and cheap, and the
algorithm has no lower threshold. The lower bound of 1000 and the upper bound of 10¹² are
properties of the web layer's selector, not of the mathematics
([Layer ownership of bounds](/rules/layer-ownership-of-bounds.md)). Progress is emitted about
ten times a second: numbers found so far plus processed count.

# Related

[Search progress](/model/search-progress.md) ·
[Search Armstrong numbers](/usecases/search-armstrong-numbers.md) ·
[Block-wise digit-power search decision](/decisions/blockwise-digit-power-search.md)

# Citations

- `internal/z02/armstrong.go` — power tables, block search, worker pool, ticker-driven
  progress, cancellation contract, verified value list.
- `internal/z02/armstrong_test.go`, `cancel_test.go` — explicit expected value lists; the
  cancelled search returns partial state plus the cancellation error.
