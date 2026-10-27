---
type: UseCase
title: Plot a function
description: "Submit an expression and interval, receive a segmented curve with asymptote positions and an applied clamp, with requested vs used sample counts disclosed."
tags: [use-case, z06, synchronous]
timestamp: "2026-10-27T16:32:44Z"
---

# Actor · goal

Learner — plot a function they type, over an interval they choose, with vertical asymptotes
drawn correctly.

# Main flow (synchronous)

1. Submit an expression, the interval bounds, an optional sample count (default 1601).
2. Edge checks: expression ≤ 2000 characters; the interval must be strictly increasing (NaN
   bounds included in the rejection).
3. Parse the expression ([Expression](/model/expression.md)); a parse failure is a 400 with the
   offending token, character, arity, or nesting depth named.
4. Sample uniformly, derive the outlier threshold, clamp outliers to ±YClamp keeping them,
   segment the curve, and locate asymptotes ([Plot](/model/plot.md)).
5. Response: segments, asymptotes, clamp value, and `samples` alongside `samplesRequested`
   ([No silent truncation](/rules/no-silent-truncation.md)).
6. The client draws each segment as a separate path and a dashed line at every asymptote.

# Alternative flows

- **Pathological nesting** (a million parentheses) → 400, and the process survives; a stack
  overflow would kill the server
  ([Expression depth bound](/invariants/expression-depth-bound.md)).
- **Fewer than ten finite samples** → an empty plot with a default clamp, not a failure.
- **Sample count out of range** → clamped to 501…20001 and disclosed.
- **`2e3`** → 2000; **`3e`** → 3·e; an unknown character is named as a character.

# Verified outcome

`1/x` on [−5,5]: two branches, an asymptote at x ≈ 0, `yclamp` 8.36 with max drawn |y| 8.36.
`tan(x)` on [−6.3, 6.3]: four asymptotes.

# Citations

- `internal/z06/expr.go` — `PlotFun`, thresholds, segmentation, clamping.
- `cmd/wdi-server/main.go` — the plot endpoint (length bound, default samples).
- `internal/z06/clamp_test.go`, `depth_test.go`, `expr_test.go`; `.verify/gate.sh` probe.
