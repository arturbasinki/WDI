---
type: DataFlow
title: Browser client flow
description: "How the client moves between views: route change, disposal of the previous view's resources, catalog-driven naming, escaped rendering, polling, and the route race guard."
tags: [dataflow, frontend]
timestamp: "2026-10-27T16:32:44Z"
---

# Path

Fragment change → disposers of the old view → catalog lookup → view render → user input →
API call → escaped render; with timers registered and released along the way.

# Stages

1. **Route change.** A new sequence number is taken; the registered disposers of the current
   view run first (intervals, animation timers, global key handlers, running job polls).
2. **Task naming.** The view's heading and accent come from the catalog entry for the task id
   — the client holds no titles ([Task catalog](/domains/task-catalog.md)).
3. **Render.** The view template is built; every foreign value passes through escaping
   ([XSS escaping](/rules/xss-escaping.md)).
4. **Interact.** Submitting starts a synchronous call or a job; resources the view starts are
   registered as disposers, and a running job is cancelled when the view is left.
5. **Present.** Results render; a progress bar is a labelled progress indicator with a live
   value ([Accessibility](/quality/accessibility.md)).
6. **Race guard.** A slow render of an older route checks its sequence before writing; if a
   newer route won, it discards its result.

# Data that never leaves the browser

Interactive state: the partially filled sudoku board, the animation position, form values,
the accent colour, and the client-side debug state (exposed read-only through a frozen object
for automated checks).

# Related

[Browser client context](/contexts/browser-client.md) ·
[View resource disposal](/rules/view-resource-disposal.md) ·
[Browser runtime](/external-systems/browser-runtime.md)

# Citations

- `cmd/wdi-server/web/app.js` — `route`, `routeSeq`, `runDisposers`, `taskHead`, `accentOf`,
  `esc`, `poll`, `window.wdiDebug`.
- `cmd/wdi-server/web/styles.css` — tokens incl. the contrast-corrected muted text token.
