---
type: Decision
title: Block-wise digit-power search
description: "Search narcissistic numbers in blocks of ten sharing a leading part, computing the digit-power sum once per block, parallelised across cores to reach 10^12."
tags: [decision, adr, performance, z02]
timestamp: "2026-10-27T16:32:44Z"
---

# Context

The exercise asks for as many narcissistic/Armstrong numbers as possible. A straightforward
scan recomputes the sum of digit powers for every candidate, and the original Python reached
only 10⁸. The claim to be made is about reach, so cost per candidate is the deliverable.

# Decision

Iterate in **blocks of ten consecutive numbers sharing the same leading part**. For a block,
compute the digit-power sum of the leading part once; the ten candidates differ only by
`pow[lastDigit]`, precomputed in a table for every digit count 1..19. Ranges are grouped per
digit count so the exponent is constant within a block, and the range is split into 4 000 000
blocks consumed by a worker pool sized to the available cores.

# Consequences

Roughly one division, one addition and one comparison per candidate — about a tenfold
improvement over the per-number formulation, and no digit-carry recomputation at all. Reach
grows from 10⁸ to a verified 10¹² (41 numbers, 31 minutes on 12 cores), with progress reported
live. The digit tables are bounded at 19 digits because 9¹⁹ must fit the integer range, which
is a mathematical ceiling of the formulation, not a configuration.

# Cost paid

The search machinery (tables, block partitioning, pool, progress ticker) is more complex than a
naive scan and must be cancellation-safe: the pool's input channel is always closed so workers
cannot outlive a cancelled search.

# Related

[Narcissistic numbers](/domains/narcissistic-numbers.md) ·
[Cancellation granularity](/rules/cancellation-granularity.md) ·
[Search progress](/model/search-progress.md)

# Citations

- `internal/z02/armstrong.go` — `digitPow`/`pow10` tables, `chunkSearch`, `blockSize`,
  the worker pool, `maxDigits` with its integer-range justification.
- `internal/z02/armstrong_test.go` — the verified 41-value list.
