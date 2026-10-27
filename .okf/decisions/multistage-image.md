---
type: Decision
title: Two-stage reproducible image
description: "Build a two-stage, offline-reproducible container image whose runtime stage carries only the static executable, runs non-root, and declares the healthcheck in the image."
tags: [decision, adr, deployment, container]
timestamp: "2026-10-27T16:32:44Z"
---

# Context

The recommended way to run the product must require nothing on the host — no toolchain, no
module proxy, no package installation — and must be startable, healthy, stoppable and
rebuildable offline.

# Decision

Build the image in **two stages**: a compiler stage that compiles a fully static executable with
the language runtime's collection bridge disabled, path trimming and debug-info stripping, then
a minimal runtime stage that carries **only** that executable, runs as a non-root user, and
installs no packages. Pin both base images by version. Make the layer reproducible by
normalising the executable's timestamp and using link-preserving copy. Declare the healthcheck
**in the image** as well as in the compose file, deriving its URL from the container's own port
variable.

# Consequences

The runtime stage builds with no network (the server makes no outgoing requests, so it needs no
TLS trust store; the healthcheck uses the runtime's built-in fetcher). The image is ~15 MB and
identical bit-for-bit across clean rebuilds. A non-root container with a root-owned,
world-readable binary. The healthcheck is observable without compose.

The cost: two base images to keep current — both pins were bumped when their series reached
end-of-life — and the build must copy the module metadata before the sources to cache well.
Hardening beyond the user switch (read-only root filesystem, capability drops) is recorded as
an open item.

# Related

[Container platform](/external-systems/container-platform.md) ·
[Embedded content decision](/decisions/embedded-content.md) ·
[Loopback publication](/decisions/loopback-port-publication.md) ·
[Portability](/quality/portability.md) · [CI runner](/external-systems/ci-runner.md)

# Citations

- `Dockerfile` — the two stages, the static build flags, timestamp normalisation,
  link-preserving copy, `EXPOSE`/`ENV`, the image-level healthcheck.
- `docker-compose.yml` — build context, healthcheck, restart policy.
- `.dockerignore` — caches, toolchain, notes and screenshots excluded from the build context.
- `RAPORT-JAKOSC-KODU.md` §3 — the reproducibility and build-context fixes.
