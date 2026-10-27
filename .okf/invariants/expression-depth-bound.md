---
type: Invariant
title: Expression depth bound
description: "Expression parsing is bounded to 100 nesting levels and 2000 characters so pathological input is a rejection, never a process crash."
tags: [invariant, z06, safety]
timestamp: "2026-10-27T16:32:44Z"
---

# Statement

Parsing an expression descends at most 100 levels of nesting, and a submitted expression is at
most 2000 characters. Exceeding either is a clean rejection.

# Rationale

The parser recurses (expr → term → factor → power → atom → expr), so nesting depth is unbounded
recursion. One million nested `(` is a stack overflow, and a stack overflow is a fatal process
condition that the transport's own recovery cannot catch: it kills the whole server, taking the
other tasks and the running jobs with it. The depth limit turns a crash vector into a 400.

# Enforcement

The depth counter is incremented on entering the innermost parse rule and decremented on
leaving it; the length bound is enforced at the transport edge, before parsing begins.

# Verified

A one-million-parenthesis expression is answered with a 400 and the process stays alive — that
probe is part of the [Quality gate workflow](/workflows/quality-gate.md).

# Related

[Expression](/model/expression.md) · [Function plotting](/domains/function-plotting.md) ·
[Reliability](/quality/reliability.md)

# Citations

- `internal/z06/expr.go` — `maxParseDepth`, the depth guard in `atom`.
- `internal/z06/depth_test.go` — deep nesting rejected.
- `cmd/wdi-server/main.go` — the 2000-character expression bound.
- `.verify/gate.sh` — the paren-storm probe asserting the server is still alive.
