---
type: QualityAttribute
title: Reproducibility of results
description: "Reproducibility: every documented numeric result is pinned by explicit-value tests and re-verified against a live process, with determinism engineered into seeded and ranked output."
tags: [quality-attribute, iso25010, reliability, testing]
timestamp: "2026-10-27T16:32:44Z"
---

# Attribute class (ISO/IEC 25010)

**Maintainability** (functional suitability in practice) — every numeric claim the project
makes can be reproduced by running the product.

# Target

The documented results are not aspirational. Each is reproducible from a clean checkout, and
the reproduction is asserted, not asserted-about.

# Pinned results

| Claim | Value |
|---|---|
| z01 exercise example | x = 18, y = −58, z = 35 |
| z02 verified list | 41 numbers up to 10¹² (31 min on 12 cores); ≤10⁹ also unit-tested |
| z03 prime counts | 168 / 1229 / 9592 / 78 498 / 664 579 / 5 761 455 |
| z04 top words | *rzekł* 155, *tylko* 149, *hrabia* 129 |
| z05 key example | 312063 → "fpop whsbsilhoee"; proverb broken with 142807 (253 918 keys) |
| z06 plot shape | `1/x` on [−5,5]: clamp 8.36, max drawn \|y\| 8.36; `tan(x)`: 4 asymptotes |
| z07 solution counts | 8 / 1 / 476 for the three documented puzzles |
| z08 blanks | 55–58 for seeds 1–8 and 42, uniqueness proven by counting |

# Mechanisms

- Tests state **explicit, complete expected lists**, not properties, so a drifted value fails.
- Determinism is engineered where randomness is involved: seeded generation, ordered output
  from unordered collections, tie-broken rankings
  ([Test determinism](/quality/test-determinism.md)).
- Claims are re-checked against a **live** process by the gate
  ([Quality gate workflow](/workflows/quality-gate.md)).
- Cross-verification against the original implementation settled the numeric baseline.
- The honesty rule keeps reported values equal to used values
  ([No silent truncation](/rules/no-silent-truncation.md)).

# Known limit

The 10¹² sweep is a recorded measurement, not a test — rerunning it costs half an hour
([Test cost budget](/decisions/test-cost-budget.md)).

# Citations

- `README.md` — the verified-results section.
- `internal/z01..z08` tests — the explicit expected values.
- `cmd/wdi-server/main.go` — the verified narcissistic-number list served by the API.
- `.verify/gate.sh` — the live reproduction probes.
