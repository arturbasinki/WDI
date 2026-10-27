---
type: Decision
title: Task catalog as the single source
description: Make the serving-layer task registry the only source of task titles and blurbs; the frontend keeps only a route-to-view mapping.
tags: [decision, adr, catalog, frontend]
timestamp: "2026-10-27T16:32:44Z"
---

# Context

The eight tasks were listed twice: a registry in the serving layer and a copy in the frontend,
each carrying titles and blurbs. Two writers of the same fact, with nothing checking agreement.
The same duplication made a blurb change a two-file edit that was easy to half-do.

# Decision

The serving-layer registry is the **only** writer of task identity, titles and blurbs, served
over the catalog endpoint. The frontend keeps only a route table (task id → view) plus local
presentation concerns (an accent colour per task id). Headings, tiles and page context read
from the fetched catalog.

# Consequences

Editing a blurb is one edit, and the health probe, the tiles, the headings and the CI probe all
see the same text. A test asserts the registry's shape (eight entries, order, unique ids,
non-empty fields) and the no-overclaim wording, so the single source is also guarded
([Task catalog consistency](/invariants/task-catalog-consistency.md)).

The cost: the client cannot name a task before the catalog arrives, so a catalog failure is a
page-level failure rather than a degraded render; and a route render is guarded by a sequence
number so a slow catalog cannot overwrite a newer route.

# Related

[Task catalog](/domains/task-catalog.md) · [Task descriptor](/model/task-descriptor.md) ·
[Task catalog API](/contracts/task-catalog-api.md) ·
[Browser client flow](/dataflows/browser-client-flow.md)

# Citations

- `cmd/wdi-server/main.go` — the registry and its endpoint.
- `cmd/wdi-server/web/app.js` — `loadCatalog`, `catalogById`, `taskTitle`, `taskHead`,
  `ACCENTS` (the only task-keyed data kept client-side).
- `RAPORT-JAKOSC-KODU.md` §5 P3 item 15 — the duplication this closed.
