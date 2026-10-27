---
type: QualityAttribute
title: Resource boundedness
description: "Performance efficiency: every resource has an explicit bound (input sizes, narration cells, jobs, slots, timeouts), and the reach claims are measured, not assumed."
tags: [quality-attribute, iso25010, performance-efficiency]
timestamp: "2026-10-27T16:32:44Z"
---

# Attribute class (ISO/IEC 25010)

**Performance efficiency** — resource utilization and capacity, under the project's capacity
model: one process, two computations at a time, ephemeral state.

# Target

Resource use is bounded by explicit budgets, not by optimism: memory per request, per job, per
running computation, and per unit of time.

# The budget ledger

| Resource | Bound |
|---|---|
| Request body | 4 MB |
| Expression length / nesting depth | 2000 characters / 100 levels |
| Linear system size | N ≤ 100 |
| Elimination narration | 400 000 matrix cells |
| Sieve range | n ≤ 10⁸ |
| Narcissistic search reach | 10¹² at the edge, 10¹³ in the library |
| Cipher key length / keys tried | ≤ 9 / ≤ 10⁸, default 2·10⁶ |
| Plot samples | 501 … 20 001 |
| Returned cryptarithm solutions | ≤ 500 (default 50) |
| Concurrent computations | 2 |
| Retained jobs | 200, for 15 minutes after finishing |
| Shutdown grace | ≤ 5 s |

# Measured efficiency, as a deliverable

Reach is a claim of the exercises themselves: 10¹² for the number search (block-wise digit-power
search), 10⁸ for the sieve with three structures compared, sub-second for a full cryptarithm
permutation search, 1–7 ms for a sudoku puzzle, tens of milliseconds for a 100×100 linear
system. Timings the learner sees are produced under realistic conditions
([Measurement flow](/dataflows/measurement-flow.md)).

# Related

[Concurrency cap](/rules/concurrency-cap.md) · [Job retention](/rules/job-retention.md) ·
[Animation step budget](/rules/animation-step-budget.md) ·
[Reliability](/quality/reliability.md) · [Cancellation granularity](/rules/cancellation-granularity.md)

# Citations

- `cmd/wdi-server/main.go` — the bounds applied at the edge, the semaphore, the retention
  constants, the server timeouts (read 30 s, write 2 min, idle 120 s).
- `internal/z01/solver.go`, `internal/z02/armstrong.go`, `internal/z03/sieve.go`,
  `internal/z05/gronsfeld.go`, `internal/z06/expr.go`, `internal/z07/cryptarithm.go` —
  the algorithmic bounds.
- `README.md` — the measured times behind each reach claim.
