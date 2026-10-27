---
type: DataFlow
title: System composition flow
description: "How the artifact is assembled from embedded content and how data moves between browser, serving layer, solvers, jobs and the container platform at run time."
tags: [dataflow, c4, deployment]
timestamp: "2026-10-27T16:32:44Z"
---

# Shape

One deployed artifact assembled at build time, three cooperating contexts at run time.

# Build time

Source code, the single-page frontend, the canonical literary text, and the English
dictionary are compiled into one self-contained executable. Nothing is copied, installed or
configured at run time; the container's runtime stage contains that one file and nothing else.

# Run time

```
Browser ──GET static──▶ HTTP API context ──embedded FS──▶ frontend files
Learner ──POST /api/…──▶ HTTP API context
                            │ validate params (400 here, job never created)
                            ├─ sync solver call ─────────▶ [Solver library context] ──▶ JSON result
                            └─ async: acquire slot (429 if none) → create Job →
                                        goroutine runs solver with cancel + progress ──▶ Job
Learner ──GET /api/…/{id}──▶ job snapshot (status, progress, result)
Learner ──DELETE /api/…/{id}──▶ cancel active (202) / conflict if finished (409)
Container ──SIGTERM──▶ cancel all jobs → close listener within 5 s → exit 0
```

# Data at rest

Only in memory, and only for the retention window: job snapshots. Everything else is
content embedded in the artifact or computed on demand.

# Related

[System](/system.md) · [HTTP API context](/contexts/http-api.md) ·
[Solver library context](/contexts/solver-library.md) ·
[Browser client context](/contexts/browser-client.md) ·
[Embedded content decision](/decisions/embedded-content.md)

# Citations

- `cmd/wdi-server/main.go` — embedded filesystem served for `/`, API mounted on `/api/`,
  job start helper, shutdown wiring.
- `internal/z04/words.go`, `internal/z05/gronsfeld.go` — embedded corpus and dictionary.
- `Dockerfile` — compile stage → runtime stage carrying only the executable.
