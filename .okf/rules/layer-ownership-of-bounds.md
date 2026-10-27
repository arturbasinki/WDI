---
type: BusinessRule
title: Layer ownership of bounds
description: "User-facing bounds live in the web layer; the solver library enforces only mathematical limits, so small and large inputs are both honest there."
tags: [business-rule, architecture, api]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

A bound that exists to protect the *user experience* belongs to the web layer; a solver library
imposes only the bounds that are properties of the *mathematics*.

# The canonical case

The web layer accepts a narcissistic-number search limit of 1000…10¹² — exactly the range the
UI selector offers, and 10¹² is the reach that was verified. The library imposes no lower
threshold at all: searching [0,100] is correct and cheap, and the library's own ceiling
(10¹³) comes from integer range, not from taste. The library's behaviour for a small limit is
pinned by a test that asserts it answers honestly.

# Corollaries

- The web layer's bound is a policy the library can be tested against independently.
- Putting the lower bound inside the library hides it from the API contract and makes the
  library untestable at small scale.
- Where the web layer clamps instead of rejecting, it must say so
  ([No silent truncation](/rules/no-silent-truncation.md)).

# Same rule elsewhere

The number of returned cryptarithm solutions is the web layer's concern (default 50, ceiling
500, clamped); the library reads a non-positive limit as "no limit". The default word count is
applied by the caller, and the library errors on a non-positive request
([Explicit defaults](/rules/explicit-defaults.md)).

# Related

[Synchronous validation](/rules/synchronous-validation.md) ·
[Solver library context](/contexts/solver-library.md) ·
[Solution capping](/rules/solution-capping.md)

# Citations

- `internal/z02/armstrong.go` — the accepted-range comment and the `limit < 0` / ceiling checks
  only; `TestSmallLimits` in `internal/z02/armstrong_test.go`.
- `cmd/wdi-server/main.go` — the 1000…10¹² check, the default-limit application, the
  solution-limit clamp for z07.
- `internal/z07/cryptarithm.go` — `limit <= 0` means unlimited.
