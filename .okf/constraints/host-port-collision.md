---
type: Constraint
title: Host port 8080 collision
description: "Port 8080 is taken on the host by an unrelated service, so local access, probes and documentation use 127.0.0.1:8099 while the container keeps 8080."
tags: [constraint, environment, deployment]
timestamp: "2026-10-27T16:32:44Z"
---

# Constraint

Port 8080 on the development host is occupied by an unrelated service, so the product's default
port is not the port it is reached on locally.

# Consequences

- The container maps its 8080 to host 8099, loopback-only
  ([Loopback publication](/decisions/loopback-port-publication.md)); the container's own
  `PORT` and the image's healthcheck stay at 8080.
- The local verification convention is fixed: live probes run against 8099. Probing 8080
  reaches someone else's service and yields responses that look like authentication failures —
  a misleading, not a loud, failure.
- Documented access is `http://127.0.0.1:8099`, and the README explains the IPv4-only binding,
  because `localhost` may resolve to the IPv6 loopback and be refused.
- `docker compose up` on a machine where the collision is not accounted for fails with a
  port-allocation error; the documented remedy is to change the left side of the mapping.

# Nature of the constraint

Environmental, not architectural. Nothing in the product knows about 8099; the port comes from
the environment ([HTTP API context](/contexts/http-api.md)).

# Citations

- `docker-compose.yml` — the mapping and the explanatory comments.
- `README.md` — the run instructions and the port-change recipe.
- `.verify/gate.sh` — the probes against 8099.
- `RAPORT-JAKOSC-KODU.md` §7 — the occupying service and the observed 401 from probing 8080.
