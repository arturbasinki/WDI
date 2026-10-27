---
type: QualityAttribute
title: Security posture
description: "Security posture: no XSS, no command execution, no outbound calls, no traversal, bounded resource abuse, non-root and loopback-only deployment, with accepted gaps recorded."
tags: [quality-attribute, iso25010, security]
timestamp: "2026-10-27T16:32:44Z"
---

# Attribute class (ISO/IEC 25010)

**Security** — confidentiality, integrity, and the absence of injection and resource-abuse
vectors in a locally-served teaching application.

# Posture

| Vector | State |
|---|---|
| Script injection into the page | None: every externally-sourced value is escaped
  ([XSS escaping](/rules/xss-escaping.md)) — verified exhaustively across all eight views |
| Server-side command execution | None: no commands are built or executed |
| Outbound requests from the server | None: the server never calls out, so no trust store is needed |
| Path traversal of served content | Not reachable: served content is fixed at build time |
| Resource abuse by a single client | Bounded: body limit, parameter bounds, two-computation cap,
  narration budget, recursion depth ([Resource boundedness](/quality/resource-boundedness.md)) |
| Privilege level of the deployed process | Non-root user, root-owned world-readable executable |
| Network exposure | Loopback interface only ([Loopback publication](/decisions/loopback-port-publication.md)) |
| Unvalidated input reaching a solver | None: all edge validation is synchronous |

# Accepted gaps (recorded, not hidden)

No content-security policy and no cache headers; an over-long request body is answered with 400
rather than the semantically correct 413; container hardening beyond the user switch
(read-only root filesystem, capability drops, no-new-privileges) is not applied. All are
acceptable for a loopback, single-purpose, local application, and all are listed as open items.

# Related

[Static asset delivery](/contracts/static-assets.md) ·
[Error communication](/rules/error-communication.md) ·
[Container platform](/external-systems/container-platform.md) ·
[Reliability](/quality/reliability.md)

# Citations

- `RAPORT-JAKOSC-KODU.md` §4 — the exhaustive XSS verification, the accepted gaps, and the
  container privilege review.
- `cmd/wdi-server/main.go` — the body limit and the panic containment.
- `Dockerfile`, `docker-compose.yml` — user, permissions, loopback binding.
- `cmd/wdi-server/web/app.js` — the escaping helper.
