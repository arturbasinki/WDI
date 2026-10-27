---
type: Aggregate
title: Expression
description: "Aggregate of a parsed arithmetic expression over x, constants and functions: precedence, associativity, implied multiplication, arity-checked calls, bounded depth."
tags: [aggregate, z06, language]
timestamp: "2026-10-27T16:32:44Z"
---

# Type

Aggregate. A parsed function expression: a tree of operations over the variable `x`,
constants and named functions, evaluable numerically at a point.

# Grammar (the language's shape)

```
expr   := term (('+'|'-') term)*
term   := factor (('*'|'/') factor | number | ident | '(')*   # implied multiplication
factor := ('-'|'+')* power
power  := atom ('^' factor)?                                   # right-associative
atom   := NUMBER | IDENT | FUNC '(' expr (',' expr)* ')' | '(' expr ')'
```

Precedence and associativity are part of the model and are pinned: `-2^2 = -4` (unary sign
binds looser than power) and `2^3^2 = 512` (power is right-associative).

# Numbers

Mantissa with optional single decimal point, plus an exponent **only when digits follow the
sign** — so `2e3` is 2000, while `3e` and `2e+` remain 3·e. Two decimal points is an
invalid number.

# Naming

Function names are recognised as calls; unknown names are rejected with the list of available
names. Function arity is validated (`min`/`max` variadic ≥2, `pow`/`fmod` 2, everything else 1).
A rejected character in the source is reported as a whole character, not as a raw byte — an
unrecognized multi-byte character like "π" must be legible in the message.

# Invariants

- Evaluation depth is bounded ([Expression depth bound](/invariants/expression-depth-bound.md)).
- Total expression length is bounded by the serving layer (2000 characters).
- Evaluation never throws: an out-of-domain function yields a non-finite value, which the
  plotting stage treats as a gap in the curve.

# Related

[Function plotting](/domains/function-plotting.md) · [Plot](/model/plot.md) ·
[Expression depth bound](/invariants/expression-depth-bound.md)

# Citations

- `internal/z06/expr.go` — tokenizer, `tokenAliases`, the `expr/term/factor/power/atom` chain,
  the function table with `unaryArity`, `maxParseDepth`, `funcNames` in the error message.
- `internal/z06/expr_test.go`, `depth_test.go` — precedence, associativity, depth rejection.
