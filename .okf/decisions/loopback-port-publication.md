---
type: Decision
title: Loopback-only publication on host port 8099
description: "Publish the container's port on the loopback interface only and map the container's 8080 to host port 8099, because 8080 is occupied on this host."
tags: [decision, adr, deployment, security]
timestamp: "2026-10-27T16:32:44Z"
---

# Context

The container published its port without naming an interface, which binds all interfaces —
exposing a local teaching application to the LAN. Separately, port 8080 on this host belongs to
an unrelated service, so `docker compose up` failed with a port-allocation conflict
([Host port collision](/constraints/host-port-collision.md)).

# Decision

Publish the container port on the **loopback interface only**, and expose the container's
8080 under a **different host port** (8099). The container's own listening port stays 8080,
and the healthcheck derives its target from the container's port variable so the two cannot
drift.

# Consequences

Nothing is reachable outside the machine. The documented address is `http://127.0.0.1:8099`
and the README explains why `localhost` may fail: the binding is IPv4-only, and browsers
frequently resolve `localhost` to the IPv6 loopback first. Changing the host port is a one-sided
edit (the left side of the mapping), leaving the container's port and the healthcheck untouched.

The cost: an extra indirection to explain, and the local verification convention that probes
run against 8099, not 8080
([Quality gate workflow](/workflows/quality-gate.md)).

# Related

[Container platform](/external-systems/container-platform.md) ·
[Host port collision constraint](/constraints/host-port-collision.md) ·
[Graceful shutdown workflow](/workflows/graceful-shutdown.md) ·
[Security posture](/quality/security-posture.md)

# Citations

- `docker-compose.yml` — the `127.0.0.1:8099:8080` mapping, `PORT`, healthcheck via `$$PORT`,
  with the comment explaining the loopback and IPv6 reasoning.
- `Dockerfile` — `EXPOSE 8080`, `ENV PORT=8080`, the image-level healthcheck.
- `README.md` — the run instructions naming 127.0.0.1 and the port-change recipe.
- `RAPORT-JAKOSC-KODU.md` §3 — the LAN-exposure finding this closed.
