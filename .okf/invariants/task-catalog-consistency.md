---
type: Invariant
title: Task catalog consistency
description: "The task catalog has exactly eight uniquely-identified entries in order, all fields populated, no duplicated copy on the client, and no blurb that over-claims."
tags: [invariant, catalog]
timestamp: "2026-10-27T16:32:44Z"
---

# Statement

The task catalog contains exactly eight entries, in exercise order, with globally unique ids,
and with every field (id, number, title, blurb) non-empty. No client holds a second copy of a
task title, and no blurb promises more than the implementation guarantees.

# Rationale

The catalog drives the tiles, the page headings and the container health probe. A duplicate id
breaks routing, an empty title renders an empty tile, and a duplicated catalog in the client is
a guaranteed future drift: two places to edit, one of which will be forgotten.

# The over-claim clause

The blurb is a promise to a learner. The sudoku blurb must not use the word "maximum" for blank
cells, because the generator is monotonic, not optimal (55–58 blanks against a theoretical
ceiling of ~64). A regression test asserts both halves: the forbidden word, and the presence of
the honest phrases ("many empty cells", "55–58", "uniqueness").

# Related

[Task catalog](/domains/task-catalog.md) · [Task descriptor](/model/task-descriptor.md) ·
[Sudoku unique solution](/invariants/sudoku-unique-solution.md) ·
[Single-source catalog decision](/decisions/single-source-task-catalog.md)

# Citations

- `cmd/wdi-server/catalog_test.go` — eight entries, expected order, duplicate-id guard,
  non-empty field guard, `TestZ08TagBezObietnicyMaksimum`.
- `cmd/wdi-server/main.go` — the catalog registry.
- `cmd/wdi-server/web/app.js` — `loadCatalog` (client reads, never stores a title).
