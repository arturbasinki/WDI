---
type: ExternalSystem
title: Browser runtime
description: "The client-side execution environment that renders the single-page frontend, performs routing, polls jobs, and owns all UI state."
tags: [external, c4-external-system]
timestamp: "2026-10-27T16:32:44Z"
---

# What it is

The end-user's browser. It downloads the embedded static frontend once and then drives
everything through the JSON API. No build step, no framework, no bundler: the delivered
program is plain script, stylesheet and markup.

# Contract it must honour

- All data coming from the API is inserted as escaped text, never as raw markup
  ([XSS escaping rule](/rules/xss-escaping.md)).
- Timers, intervals and global handlers started by a view must be released when the route
  changes ([View resource disposal](/rules/view-resource-disposal.md)).
- A long job the user navigates away from is cancelled by the client.

# Why it is external

It executes outside the deployed artifact, holds the only copy of interactive UI state
(animations, in-progress grids), and cannot be assumed to be online for the whole
computation — which is exactly why job context must not be tied to the request
([Detached job context decision](/decisions/detached-job-context.md)).

# Citations

- `cmd/wdi-server/web/app.js`, `styles.css`, `index.html` — the frontend served verbatim.
- `cmd/wdi-server/web/app.js` — `esc`, `onLeave`/`runDisposers`, `poll`, `startJob`, `route`.
