---
type: Contract
title: Static asset delivery
description: "Delivery of the embedded single-page client (markup, stylesheet, script) as static content alongside the API."
tags: [contract, http, frontend]
timestamp: "2026-10-27T16:32:44Z"
---

# Port

Delivery of the client application itself: the frontend files, served verbatim from inside the
deployed artifact.

# Shape

- `GET /` and the static paths under it → the markup, stylesheet and script, served from the
  embedded content, with the API mounted alongside under `/api/`.
- Content is fixed at build time. Nothing downloads assets at run time, and there is no
  configuration naming a document root
  ([Embedded content decision](/decisions/embedded-content.md)).
- Routing inside the client is by URL fragment, so a deep link to a task needs no server-side
  route.

# Consequences

- Path traversal outside the embedded content is not reachable.
- The container needs no files beyond the executable, and the runtime stage builds without
  network access.
- Because the script is one unbundled file, its syntax is the checkable surface — the gate runs
  a parser check over it.
- No cache headers and no content-security policy are set; recorded as an accepted gap for a
  local teaching application, and an open item.

# Related

[Browser presentation](/domains/browser-presentation.md) ·
[Browser client context](/contexts/browser-client.md) ·
[Security posture](/quality/security-posture.md)

# Citations

- `cmd/wdi-server/main.go` — the embedded filesystem and its sub-path, the mount alongside the
  API.
- `cmd/wdi-server/web/index.html`, `app.js`, `styles.css`.
- `.github/workflows/ci.yml` — the frontend syntax check step.
- `RAPORT-JAKOSC-KODU.md` §4 — the CSP/cache gap.
