---
type: UseCase
title: Browse the task catalog
description: The learner lists the eight tasks from the server-side catalog and opens one; the client stores no task titles.
tags: [use-case, catalog]
timestamp: "2026-10-27T16:32:44Z"
---

# Actor · goal

Learner — see which exercises exist, with an honest one-line description, and open one.

# Main flow

1. The client fetches the catalog.
2. The server returns the eight entries in exercise order
   ([Task catalog contract](/contracts/task-catalog-api.md)).
3. The client renders tiles from the returned number, title and blurb, and installs the
   route table (id → view).
4. The learner picks a tile; the route changes, the previous view's resources are released
   ([View resource disposal](/rules/view-resource-disposal.md)), the new view renders.
5. The page heading is taken from the same catalog entry — the client stores no titles.

# Variants

- Catalog unavailable → the client cannot name the tasks; it shows a fallback title and no
  tiles. This is treated as a fatal page condition, not a partial render.
- A newer navigation wins a slow catalog response: a sequence guard discards the stale render.

# Acceptance

Eight tiles in order, each with a non-empty title and blurb, no duplicate target, and the
blurb for sudoku not promising a maximum
([Task catalog consistency](/invariants/task-catalog-consistency.md)). The catalog endpoint is
also the container's health probe, so this use case is what proves the process is serving.

# Citations

- `cmd/wdi-server/main.go` — the catalog handler; `cmd/wdi-server/catalog_test.go`.
- `cmd/wdi-server/web/app.js` — `loadCatalog`, `renderHome`, `taskHead`, `routeSeq`.
- `docker-compose.yml`, `Dockerfile` — healthcheck on this endpoint.
