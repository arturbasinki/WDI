---
type: Decision
title: Test cost budget
description: "Budget test cost: heavy numeric cases are short-mode-skipped, HTTP tests stub the heavy computations, and timing tests use an injected clock, keeping the race-enabled suite in seconds."
tags: [decision, adr, testing]
timestamp: "2026-10-27T16:32:44Z"
---

# Context

The verification suite ran for minutes at full CPU, dominated by the two genuinely expensive
cases: a narcissistic-number search to 10⁹ and a sieve benchmark to 10⁷. A slow suite is a
suite that gets skipped, and the gate is required to run on every change
([Quality gate workflow](/workflows/quality-gate.md)).

# Decision

Budget the suite by cost: the heaviest cases declare themselves short-mode-skipped, so the fast
set runs in about a second and the full set in seconds. The serving layer is in scope on the
same footing as the solvers, with the heavy computations replaced by deterministic stubs
([Solver capability port](/contracts/solver-package-port.md)), and no test sleeps beyond 200 ms.
Retention timing is tested with an injected clock, never by waiting.

# Consequences

Full suite with the race detector: seconds, not minutes, so it runs on every change and in CI.
The expensive verification of the headline claim (the sweep to 10¹²) is deliberately **not** in
the suite — it is a one-off, recorded measurement whose result is pinned as data
([Reproducibility](/quality/reproducibility.md)); the suite covers values ≤ 10⁹.

The cost: the headline claim is protected by a recorded list, not by a rerun, so a regression
that only appears above 10⁹ would not be caught by tests.

# Related

[Quality gate workflow](/workflows/quality-gate.md) ·
[Hosted CI runner](/external-systems/ci-runner.md) ·
[Test determinism](/quality/test-determinism.md)

# Citations

- `internal/z02/armstrong_test.go`, `internal/z03/sieve_test.go` — short-mode guards on the
  heavy ranges.
- `cmd/wdi-server/handlers_test.go` — stubbed computations, injected clock, no long sleeps.
- `.github/workflows/ci.yml` — the full race-enabled suite as a gate step.
- `README.md` — the documented durations of the full and short runs.
