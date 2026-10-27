---
type: Domain
title: Task catalog (domain)
description: "The single authoritative registry of the eight tasks (identity, number, title, blurb) that the whole presentation layer reads."
tags: [domain, supporting, catalog]
timestamp: "2026-10-27T16:32:44Z"
---

# Kind of subdomain

Supporting. It exists so the eight exercises are addressable, listed and titled from one
place, and so the UI can never disagree with the server about what the tasks are.

# Responsibility

Maintain the authoritative list of the eight tasks — identity, display number, title, and a
short truthful blurb — and serve it to the presentation layer. Nothing else holds a copy of
task titles; the frontend derives headings and tiles from this list.

# Why a single source

The catalog existed twice (server registry and a copy in the frontend) and that was recorded
as a drift hazard: editing a description in one place silently desynchronised the other.
The registry in the serving layer is now the only writer, and the client maps route → view
without naming the task.

# Blurb honesty

The blurb is a promise to the learner, so it must not exceed what the generator guarantees.
The sudoku blurb states "many empty cells (typically 55–58), guaranteed uniqueness" rather
than "maximum number of empty cells", because the digging procedure is monotonic, not optimal.

# Related

[Task descriptor](/model/task-descriptor.md) ·
[Task catalog invariants](/invariants/task-catalog-consistency.md) ·
[Task catalog contract](/contracts/task-catalog-api.md) ·
[Single-source catalog decision](/decisions/single-source-task-catalog.md) ·
[System](/system.md)

# Citations

- `cmd/wdi-server/main.go` — the task registry type and the eight-entry catalog; the
  catalog endpoint handler.
- `cmd/wdi-server/catalog_test.go` — eight entries, unique ids, no empty fields, and the
  no-overclaim guard for the sudoku blurb.
- `cmd/wdi-server/web/app.js` — `loadCatalog`, `catalogById`, `taskTitle`, `taskHead`:
  titles read from the API only.
