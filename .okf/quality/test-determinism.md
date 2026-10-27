---
type: QualityAttribute
title: Test determinism
description: "Test determinism: explicit full expected values, seeded and ordered output, injected clocks, stubbed heavy computations, and job goroutines drained between tests."
tags: [quality-attribute, testing, iso25010]
timestamp: "2026-10-27T16:32:44Z"
---

# Attribute class (ISO/IEC 25010)

**Maintainability** — testability: a test that passes twice passes for the same reason.

# Rules of test writing in this project

1. **Explicit, complete expected values.** A test lists the expected numbers, words, solutions
   or rows in full, rather than asserting a property of them.
2. **Deterministic output from non-deterministic machinery.** Map iteration is ordered before
   return; rankings break ties alphabetically; sudoku generation is driven by an explicit seed;
   the cipher key enumeration is lexicographic with zero padding.
3. **No clock dependence.** Retention is tested with an injected clock
   ([Lazy retention](/decisions/lazy-retention-no-timers.md)); no test sleeps beyond 200 ms.
4. **No real heavy computation in the HTTP layer.** The three heavy capabilities are replaced by
   deterministic stubs through the port's replaceable call points, and the harness drains job
   goroutines before restoring the real ones
   ([Solver capability port](/contracts/solver-package-port.md)).
5. **Cost is budgeted.** The heaviest numeric cases are skipped in short mode
   ([Test cost budget](/decisions/test-cost-budget.md)).
6. **One test per acceptance criterion**, phrased in the language of the criterion, asserting
   the status code **and** the message text
   ([Background job policy](/rules/background-job-policy.md)).
7. **Red first.** A reported defect becomes a failing test before it becomes a fix
   ([Quality gate workflow](/workflows/quality-gate.md)).
8. **Concurrency is in scope** — the suite runs under the race detector
   ([Race freedom](/quality/race-freedom.md)).
9. **Ordering of results is part of the contract** — the catalog order, the sieve structure
   order, the ascending prime lists, the sorted found-numbers list.

# Citations

- `cmd/wdi-server/handlers_test.go` — the fake clock, the seams, the drain-before-restore rule.
- `internal/z02/armstrong_test.go`, `internal/z03/sieve_test.go` — explicit lists, short-mode
  guards; `internal/z03/order_test.go` — ordered map output.
- `internal/z08/seed_test.go` — seeded determinism.
- `cmd/wdi-server/catalog_test.go` — catalog order and content.
