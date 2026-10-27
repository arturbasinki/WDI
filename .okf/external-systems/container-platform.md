---
type: ExternalSystem
title: Container platform
description: "The container engine and base images used to build and run the product; recommended deployment target, not a code dependency."
tags: [external, deployment]
timestamp: "2026-10-27T16:32:44Z"
---

# Role

The container engine plus two pinned base images (a compiler stage and a minimal Linux
runtime stage) build and run the product. Deployment via container is the documented
recommended path because the host then needs no toolchain at all.

# Interfaces with the system

- Build: source is copied in, the compile stage produces one static executable, the runtime
  stage copies only that file. No package installation at runtime, no network in the build.
- Runtime: the container health is observed by polling the task-catalog endpoint
  ([Task catalog contract](/contracts/task-catalog-api.md)); the process runs as a non-root user.
- Shutdown: the platform sends a termination signal and expects the process to stop the
  running computations and close within the grace window
  ([Graceful shutdown workflow](/workflows/graceful-shutdown.md)).

# Notable coupling to the host

The published port is bound to the loopback interface only, and the host port differs from
the container port because port 8080 on this machine belongs to an unrelated service —
see [Host port collision constraint](/constraints/host-port-collision.md) and
[Loopback-only publication decision](/decisions/loopback-port-publication.md).

# Citations

- `Dockerfile` — two stages, static build, non-root user, healthcheck on the task list.
- `docker-compose.yml` — loopback-only port map, restart policy, healthcheck, `PORT`.
- `.dockerignore` — excludes caches, toolchain, notes, screenshots from build context.
