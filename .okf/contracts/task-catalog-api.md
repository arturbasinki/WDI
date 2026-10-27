---
type: Contract
title: Task catalog API
description: "Read-only contract exposing the eight task descriptors, consumed by the client, the container healthcheck and the gate probe."
tags: [contract, api, catalog]
timestamp: "2026-10-27T16:32:44Z"
---

# Port

A read-only listing of the eight tasks: the single source of task identity, titles and blurbs
for every client, and the health probe of the deployment.

# Shape

`GET /api/tasks` → 200

```json
{ "tasks": [ { "id": "z01", "num": "01", "title": "…", "tag": "…" }, … ] }
```

Exactly eight entries, in exercise order, unique ids, every field non-empty
([Task catalog consistency](/invariants/task-catalog-consistency.md)). No parameters, no
authentication, no caching semantics declared.

# Consumers

The client (tiles and page headings), the container healthcheck (image and compose), the CI/
gate probe (`tasks=8`), and any operator.

# Contract strength

Because there is no second copy, changing a blurb in the registry is the whole change. A test
asserts the count, order, uniqueness, non-emptiness, and the no-overclaim wording for sudoku.

# Related

[Task catalog](/domains/task-catalog.md) · [Task descriptor](/model/task-descriptor.md) ·
[Browse the task catalog](/usecases/browse-task-catalog.md)

# Citations

- `cmd/wdi-server/main.go` — the catalog handler and registry.
- `cmd/wdi-server/catalog_test.go` — the contract assertions.
- `Dockerfile`, `docker-compose.yml`, `.verify/gate.sh` — the health probe on this endpoint.
