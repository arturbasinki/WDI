---
type: Actor
title: Learner
description: "A person working through the exercises: supplies inputs, watches results and progress, cancels work they no longer want."
tags: [actor, c4-person]
timestamp: "2026-10-27T16:32:44Z"
---

# Role

Primary human actor. Opens the start screen, picks a task, enters task data (a system of
equations, a search limit, a cipher text and key, a function formula, a sudoku grid), and
reads the answer. Starts long searches, watches their progress, and leaves the page —
which is expected to stop the work, not to orphan it.

# Goals

- See the answer for the canonical example of each task without preparing any file.
- Supply their own data and get an honest answer, including when part of their input was
  dropped, clamped, or capped — see [No silent truncation](/rules/no-silent-truncation.md).
- Understand how the answer was produced (the elimination steps animation, the generator
  explanation panel, the benchmark bars).

# Consequences for the system

Every capability needs a bounded, documented input range and a human-readable rejection
when the input is outside it. Progress must be pollable. The whole interface is in Polish.

# Related

[Browse the task catalog](/usecases/browse-task-catalog.md) ·
[Browser client context](/contexts/browser-client.md)

# Citations

- `cmd/wdi-server/web/app.js` — the eight views, hash routing, progress pollers.
- `cmd/wdi-server/web/index.html` — shell, header, footer description of the product.
