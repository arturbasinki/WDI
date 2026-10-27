---
type: Decision
title: "Standard-library-only, single binary"
description: "Implement all eight tasks and the serving layer on the standard library alone, shipped as one self-contained executable, replacing every third-party library with first-party code."
tags: [decision, adr, dependencies]
timestamp: "2026-10-27T16:32:44Z"
---

# Context

Eight algorithm exercises, to be presented as one interactive web page, on a machine where
external dependencies are a liability: the working environment has no access to a module proxy,
and the deliverable must run anywhere.

# Decision

Implement everything on the language's **standard library only**, with zero external
dependencies, and ship it as **one executable** in which the frontend and the reference content
are compiled in. Each task is a separate internal module; the serving layer is a single process
that mounts the API and the static client side by side.

The libraries the original Python solutions leaned on are replaced by first-party
implementations: a computer-algebra dependency by a hand-written recursive-descent parser, a
plotting library by sampled segments plus hand-drawn charts, an array library by a packed
bitmap, a counter by an explicit tally and sort, a progress-bar library by polled progress.

# Consequences

**Positive.** Nothing to install; the container's runtime stage is the executable and nothing
else, and builds without network. Reproducible results across machines. A rewrite in another
stack has no third-party surface to replicate. Module listing confirms a single module.

**Negative.** The expression language, the chart rendering and the progress machinery are
hand-maintained: the expression language supports a fixed function set, charts are limited to
what the views draw, and progress requires a cadence rule that a library would have given.

**Confirmed by.** `go list -m all` shows one module; the image is ~15 MB; the gate builds the
binary and the image on every run.

# Related

[System](/system.md) · [Embedded content decision](/decisions/embedded-content.md) ·
[Standard-library-only constraint](/constraints/stdlib-only.md) ·
[Portability](/quality/portability.md)

# Citations

- `go.mod` — a module declaration with no requirements.
- `cmd/wdi-server/main.go` — the embedded frontend and the API mount.
- `internal/z06/expr.go`, `internal/z03/sieve.go`, `internal/z04/words.go` — the
  replacements for the original libraries.
- `README.md`, `RAPORT-JAKOSC-KODU.md` §6 — the zero-dependency claim and its verification.
