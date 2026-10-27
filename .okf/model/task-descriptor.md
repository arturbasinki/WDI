---
type: Entity
title: Task descriptor
description: "Entity for one catalogued task: stable id (the exercise number), display number, title and blurb, served as the single source of task identity."
tags: [entity, catalog]
timestamp: "2026-10-27T16:32:44Z"
---

# Type and identity

Entity. One entry of the task catalog, identified by its task id (`z01`…`z08`), which is also
the URL segment of that task's endpoints and the route key in the client. The id is stable and
meaningful: it is the exercise number.

# Schema

- `id` — the task key (unique across the catalog),
- `num` — the two-digit display number,
- `title` — the display title,
- `tag` — a one-line blurb shown on the tile and in the page header.

# Behaviour

Immutable at run time; the catalog is read-only, served by one endpoint, and is the only
source of task titles and blurbs. The client resolves id → view locally and never stores a
title.

# Ordering

Catalog order is the exercise order and is asserted: the catalog returns the eight entries in
the expected sequence. The sudoku blurb is additionally constrained so it cannot promise a
maximum ([Task catalog invariants](/invariants/task-catalog-consistency.md)).

# Related

[Task catalog](/domains/task-catalog.md) ·
[Task catalog contract](/contracts/task-catalog-api.md) ·
[Single-source catalog decision](/decisions/single-source-task-catalog.md)

# Citations

- `cmd/wdi-server/main.go` — the task definition type with its serialised field names, the
  eight-entry catalog, the catalog handler.
- `cmd/wdi-server/catalog_test.go` — eight entries, order, uniqueness, non-empty fields,
  no-overclaim blurb rule.
- `cmd/wdi-server/web/app.js` — `loadCatalog`, `catalogById`, `taskHead`.
