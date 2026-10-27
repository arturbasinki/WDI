---
type: Domain
title: Browser presentation (domain)
description: "The browser-side presentation domain: one hash-routed single-page client rendering all eight task views without a build step."
tags: [domain, supporting, frontend]
timestamp: "2026-10-27T16:32:44Z"
---

# Kind of subdomain

Supporting, generic. The browser-side presentation of the eight tasks: a hash-routed
single-page application with hand-written views, SVG charts, animations and progress bars.

# Responsibility

Render the start screen and the eight task views from catalog data, collect inputs, call
the API, present results honestly, keep interactive resources alive only as long as their
view is on screen, and meet WCAG 2.1 AA for the text it renders.

# Model it maintains

- One view definition per task; a route table maps a task id to its view.
- A per-route resource registry: intervals, animation timers and global key handlers are
  registered when a view starts and released when the route changes
  ([View resource disposal](/rules/view-resource-disposal.md)).
- A job poller: start, poll progress, stop on completion, on error (with a visible banner),
  and on leaving the page.
- A small, read-only, frozen debug surface for automated checks — the UI remains the writer.

# Deliberate constraints

No framework, no build step, no runtime dependencies: the whole client is delivered as
plain files embedded in the deployed artifact. Charts are hand-drawn SVG. All text is Polish.

# Related

[Browser client context](/contexts/browser-client.md) ·
[Browser runtime](/external-systems/browser-runtime.md) ·
[Accessibility quality attribute](/quality/accessibility.md) ·
[XSS escaping rule](/rules/xss-escaping.md) ·
[Static asset contract](/contracts/static-assets.md)

# Citations

- `cmd/wdi-server/web/app.js` — `VIEWS`, `route`/`routeSeq`, `onLeave`/`runDisposers`,
  `poll`, `startJob`, `esc`, `fmtSec`, `accentOf`, `window.wdiDebug`.
- `cmd/wdi-server/web/styles.css` — design tokens incl. the muted text token that had to be
  darkened for contrast.
- `.verify/contrast.py` — token-pair contrast computation used as a gate step.
