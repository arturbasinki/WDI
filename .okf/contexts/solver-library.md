---
type: BoundedContext
title: Solver library context
description: "The bounded context of pure computational solvers: one module per task, no transport, no shared state, errors for invalid input, cooperative cancellation."
tags: [bounded-context, hexagonal, port-side]
timestamp: "2026-10-27T16:32:44Z"
---

# Boundary

The library of eight solvers: pure computation, no transport, no I/O, no notion of a request,
a job, a port, or a screen. One task per module; modules do not know about each other.

# Model consistency inside the boundary

- Every capability is a function of its inputs, optionally taking a cooperative cancellation
  signal and a progress callback.
- Progress callbacks are never required to be non-nil; a nil callback means "no progress".
- Invalid input is an error, never a substituted default. A lower bound on a search limit,
  a default top-N, a default solution limit are all the caller's business
  ([Layer ownership of bounds](/rules/layer-ownership-of-bounds.md),
  [Explicit defaults](/rules/explicit-defaults.md)).
- A cancelled computation returns the partial state **together with** the cancellation —
  it never reports success.
- Nothing panics across the boundary: the serving layer converts a panic into a job error.
- Reference data (canonical text, dictionary) is content of the module, delivered with it.

# Language

Polish. Package doc comments state the exercise, the approach, and the boundaries.

# Provided to

[HTTP API context](/contexts/http-api.md) through the
[Solver package port](/contracts/solver-package-port.md).

# Citations

- `internal/z01/solver.go` … `internal/z08/sudoku.go` — one module per task.
- `internal/z02/armstrong.go`, `internal/z03/sieve.go`, `internal/z05/gronsfeld.go` —
  `ctx` + `onTick` signatures, nil-callback guards, partial-state-on-cancel contract.
