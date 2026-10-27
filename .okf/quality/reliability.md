---
type: QualityAttribute
title: Reliability and crash containment
description: "Reliability: no documented request can crash the process or exhaust memory — bounded inputs, narration, recursion, concurrency, crash containment, overflow guards, and live probes."
tags: [quality-attribute, iso25010, reliability, safety]
timestamp: "2026-10-27T16:32:44Z"
---

# Attribute class (ISO/IEC 25010)

**Reliability** — maturity and fault tolerance: the product must not be taken down by a
legitimate request.

# Target

No input accepted by the documented API can crash the process, exhaust its memory, or wedge a
worker beyond the retention and concurrency bounds.

# Mechanisms

- **Bounded inputs at the edge** — every parameter range, and the request body size
  ([Synchronous validation](/rules/synchronous-validation.md)).
- **Bounded narration** — the elimination recording's cell budget, closing the largest
  known memory blow-up (a 100×100 system that cost ~4.7 GB now costs ~1.18 MB)
  ([Animation step budget](/rules/animation-step-budget.md)).
- **Bounded recursion** — parse depth capped, so pathological nesting is a 400 and not a stack
  overflow that the transport cannot recover from
  ([Expression depth bound](/invariants/expression-depth-bound.md)).
- **Bounded concurrency** — two computations at a time
  ([Concurrency cap](/rules/concurrency-cap.md)); bounded state
  ([Job retention](/rules/job-retention.md)).
- **Crash containment** — a crash inside a computation becomes a job `error`, because the
  worker runs outside the transport's own recovery.
- **Overflow guards** — a product too large for the integer range is excluded from matching
  rather than compared; digit-power tables are capped at the range limit.
- **Termination guarantees** — a worker pool's input is always closed, so cancellation cannot
  leave a search running forever.

# Probes

The gate asserts, against a live process: an out-of-range limit → 400; the largest linear
system within a size budget; a million nested parentheses → 400 with the process still alive;
the catalog still served afterwards; SIGTERM → exit 0.

# Citations

- `.verify/gate.sh` — the crash-vector probes.
- `cmd/wdi-server/main.go` — body limit, panic recovery, semaphore, retention, shutdown.
- `internal/z01/solver.go`, `internal/z06/expr.go`, `internal/z07/cryptarithm.go`,
  `internal/z02/armstrong.go` — the guards.
- `RAPORT-JAKOSC-KODU.md` §3, §5 — the crash vectors found and closed.
