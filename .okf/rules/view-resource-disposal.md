---
type: BusinessRule
title: View resource disposal
description: "Views register every timer, listener and poller they start and release them on route change; a stale route may not overwrite a newer one."
tags: [business-rule, frontend, resources]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

Anything a view allocates — intervals, animation timers, global key handlers, pollers, running
jobs — must be registered when the view starts and released when the route changes. A view that
is no longer on screen must own nothing.

# Mechanism

A per-route disposer registry. A view calls the registration helper with a release function;
on route change the registry is drained before the new view renders. Intervals use disposers,
global key handlers use abortable listeners, and the poller exposes a stop that fires on
completion, on error (with a visible banner), and on leaving the page.

# Defects this rule exists to prevent

- An animation interval from the linear-systems view kept running after navigation.
- A sudoku view assigning a global key handler, overwriting any other and leaking it.
- A poller swallowing fetch failures: the progress bar hung forever with no message.
- A job left running after its view is gone, holding one of the two computation slots.

# Race discipline

Routing is asynchronous. A sequence number guards route completion, so a slower route that
finishes after a newer one must not overwrite it. The client's debug state object is frozen with
getter access — read-only from outside, written only by the UI — so automated checks cannot
mutate it.

# Related

[Browser presentation](/domains/browser-presentation.md) ·
[Browser client context](/contexts/browser-client.md) ·
[Cancellation semantics](/rules/cancellation-semantics.md)

# Citations

- `cmd/wdi-server/web/app.js` — `disposers`, `onLeave`, `runDisposers`, `poll` (with `stop`),
  `routeSeq`, `startJob`'s cancel-on-leave, `window.wdiDebug`.
- `RAPORT-JAKOSC-KODU.md` §5 P1 (frontend) — the leaked `playTimer` and `onkeydown` findings.
