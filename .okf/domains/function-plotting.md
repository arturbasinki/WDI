---
type: Domain
title: Function plotting (z06)
description: "Parse a learner-written function expression and sample it over an interval, detecting vertical asymptotes numerically and returning the curve as segments plus a clamp threshold."
tags: [domain, core, z06, plotting, parsing]
timestamp: "2026-10-27T16:32:44Z"
---

# Kind of subdomain

Core (the exercise itself). Draw the graph of a single-variable function given as a formula
and an interval, correctly handling vertical asymptotes.

# Language the learner writes

A small arithmetic expression language, parsed by a hand-written recursive-descent parser
(no external computer-algebra dependency): `+ - * / ^` (right-associative power),
parentheses, unary signs, the variable `x`, the constants `pi` and `e`, and functions
sin cos tan asin acos atan sinh cosh tanh sqrt cbrt ln log log10 exp abs floor ceil round
min max pow fmod sigmoid. Implied multiplication is accepted (`2x`, `x sin(x)`). Scientific
notation is read in numbers (`2e3` = 2000) while a bare `e` remains the constant.

# The hard part: asymptotes

Vertical asymptotes are found numerically, not symbolically. From the finite samples:
- a threshold from the distribution of |y| (median plus 15·IQR) marks outlier samples,
- a sign change between neighbouring samples with a jump above 4·IQR marks an asymptote,
- a non-finite sample next to a finite one also marks one.

The plot is returned as **segments** that asymptotes never cross, plus the asymptote
positions, so the client draws separate paths and a dashed line at each asymptote.

Outlier samples are **clamped** to ±YClamp and kept, so a branch runs to the frame instead
of ending mid-air below the threshold. The applied clamp value is returned.

# Honest sampling

Requested sample count and the count actually used (clamped to 501…20001) are returned
separately. An interval whose lower bound is not below its upper bound is rejected, including
NaN bounds. Parse depth is bounded so pathological nesting is a clean rejection and not a
crashed process ([Expression depth bound](/invariants/expression-depth-bound.md)).

# Related

[Expression](/model/expression.md) · [Plot](/model/plot.md) ·
[Plot segment](/model/plot-segment.md) · [Plot clamping](/model/plot-clamping.md) ·
[Plot a function](/usecases/plot-function.md) ·
[Numeric asymptote detection decision](/decisions/numeric-asymptote-detection.md)

# Verified results

`1/x` is drawn as two branches with an asymptote at x ≈ 0; `tan(x)` shows 4 asymptotes on
[−6.3, 6.3]; for `1/x` on [−5,5] the clamp is 8.36 and the maximum drawn |y| is 8.36, i.e. the
branch reaches the frame.

# Citations

- `internal/z06/expr.go` — tokenizer (incl. exponent notation and rune-correct error
  messages), grammar, AST evaluation, function table with arity checks, depth limit,
  percentile thresholds, segmentation and clamping.
- `internal/z06/expr_test.go`, `clamp_test.go`, `depth_test.go` — precedence/associativity
  (`-2^2 = -4`, `2^3^2 = 512`), real clamping, depth rejection.
