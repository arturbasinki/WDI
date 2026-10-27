---
type: QualityAttribute
title: Portability and single-artifact delivery
description: "Portability (ISO 25010): one self-contained artifact, no configuration, no data files, no external services, buildable offline; the bundle describes it language-agnostically."
tags: [quality-attribute, iso25010, portability]
timestamp: "2026-10-27T16:32:44Z"
---

# Attribute class (ISO/IEC 25010)

**Portability** — installability, substitutability, and the ability to be re-implemented.

# Target

The whole product travels as one self-contained executable: no configuration files, no data
directory, no external service, no database, no build step for the client. A fresh machine with
a container engine runs it in one command; a machine with the language toolchain runs it in two.

# How it is kept

- Everything the product reads is content compiled into it
  ([Content embedded](/decisions/embedded-content.md)).
- Zero dependencies make the build independent of a package supply chain
  ([Standard-library-only](/constraints/stdlib-only.md)).
- The runtime stage of the image contains only the executable and builds offline
  ([Two-stage image](/decisions/multistage-image.md)).
- The only runtime setting is the listening port, read from the environment.

# Rewrite-readiness test

The bundle's conceptual concepts ([domains](/domains/linear-systems.md), [rules](/rules/no-silent-truncation.md),
[contracts](/contracts/background-job-api.md), [invariants](/invariants/job-status-monotonic.md))
describe the product without naming a language construct: the input bounds, the status-code
policy, the progress cadences, the honesty fields and the numeric reference values are
expressible in any stack. The citations are evidence, not a translation guide.

# Known limits

The language version is pinned by the module file and the image (they must not drift); the
local tooling scripts are environment-specific
([Workspace cache](/constraints/workspace-cache.md)).

# Citations

- `README.md` — native and container run instructions, image size, zero-dependency claim.
- `go.mod`, `Dockerfile`, `docker-compose.yml` — the portable surface.
- `RAPORT-JAKOSC-KODU.md` §6 — the verified zero-dependency and image hygiene claims.
