---
type: System
title: WDI — eight algorithmic tasks as one web application
description: "A self-contained web application presenting eight academic algorithm exercises, each as an interactive graphical page over a JSON API, shipped as a single dependency-free deployable."
tags: [scope, c4-l1]
timestamp: "2026-10-27T16:32:44Z"
---

# Purpose

Present eight independently specified algorithmic exercises (the course task set) as one
coherent product: a start screen of eight task tiles, each opening an interactive, graphical
implementation of that task. Every task's answer must be reproducible from the running system,
not from a written report — numeric claims in the documentation are expected to be verifiable
live through the same API the browser uses.

# Scope

In scope:
- Eight solver capabilities: linear systems (z01), narcissistic/Armstrong numbers (z02),
  prime sieving with a three-structure benchmark (z03), word-frequency ranking (z04),
  Gronsfeld polyalphabetic cipher incl. dictionary attack (z05), function plotting with
  vertical asymptotes (z06), cryptarithm puzzles (z07), unique-solution sudoku generation
  and solving (z08). See [Task catalog](/domains/task-catalog.md).
- A presentation surface: task tiles, per-task views, results, progress indication.
- A background-computation facility for the three long-running capabilities.
  See [Background computation](/domains/background-computation.md).
- A verification discipline (quality gate) that keeps the documented claims true.
  See [Quality gate workflow](/workflows/quality-gate.md).

Out of scope: user accounts, persistence, authentication, multi-node operation, any
database. Task state is deliberately ephemeral.

# Responsibilities

1. Expose each solver as a stable, self-describing capability with explicit input bounds.
2. Keep long computations observable (progress), cancellable, and resource-bounded.
3. Be deployable as one artifact with all content (texts, dictionaries) inside it.
4. Never let a documented result drift from what the running system returns.

# Composition

Three cooperating contexts: [Solver library context](/contexts/solver-library.md) (pure
computation), [HTTP API context](/contexts/http-api.md) (transport, validation, job policy),
[Browser client context](/contexts/browser-client.md) (views and interaction).
See [System composition data flow](/dataflows/system-composition.md).

# Citations

- `README.md` — product description, run instructions, verified results.
- `cmd/wdi-server/main.go` — single entry point: HTTP, JSON API, job registry, embedded frontend.
- `internal/z01..z08` — the eight solver packages.
- `cmd/wdi-server/web/` — the frontend (one HTML, one stylesheet, one script).
