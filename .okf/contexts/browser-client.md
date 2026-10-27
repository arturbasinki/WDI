---
type: BoundedContext
title: Browser client context
description: "The bounded context of the delivered single-page client: routing, views, polling, rendering, and disciplined release of interactive resources."
tags: [bounded-context, frontend]
timestamp: "2026-10-27T16:32:44Z"
---

# Boundary

The client application delivered as static files: routing, views, input collection, progress
polling, rendering. It owns all interactive UI state and holds no computation authority —
every answer comes from the API.

# Model consistency inside the boundary

- Task titles and blurbs come from the server catalog; the client stores only a route → view
  mapping ([Task catalog domain](/domains/task-catalog.md)).
- A route change releases everything the previous view allocated
  ([View resource disposal](/rules/view-resource-disposal.md)); a slower route losing a race
  to a newer one must not overwrite it.
- A running job is cancelled when its view is left.
- All API data reaches the DOM escaped ([XSS escaping](/rules/xss-escaping.md)).
- A poller that fails must say so visibly and stop, not hang a progress bar silently.

# Contract with the browser runtime

Single-page, hash routing, no framework, no build, hand-drawn SVG, Polish copy, WCAG 2.1 AA
for the text it renders ([Accessibility](/quality/accessibility.md)).

# Related

[Browser presentation](/domains/browser-presentation.md) ·
[Browser client data flow](/dataflows/browser-client-flow.md) ·
[Static asset contract](/contracts/static-assets.md)

# Citations

- `cmd/wdi-server/web/app.js` — `VIEWS`, `route`/`routeSeq`, `onLeave`/`runDisposers`,
  `poll`, `startJob`, `esc`, chart helpers.
- `cmd/wdi-server/web/index.html`, `styles.css`.
