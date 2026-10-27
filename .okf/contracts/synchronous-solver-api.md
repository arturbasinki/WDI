---
type: Contract
title: Synchronous solver API
description: "The synchronous JSON port for the bounded-cost capabilities: request shapes, success payloads, validation and echo of applied values."
tags: [contract, api, synchronous]
timestamp: "2026-10-27T16:32:44Z"
---

# Port

Five capabilities answer inside the request, because their worst case is bounded:
linear systems, word ranking (canonical and submitted text), function plotting, cryptarithms,
sudoku generation and solving.

# Endpoints

| Endpoint | Request essentials | Success |
|---|---|---|
| `POST /api/z01/solve` | `text` (the system as typed) | 200 status, x, RREF, rank, steps, `stepsTruncated` |
| `GET /api/z04/top20` | optional query `n` (1…100) | 200 `source`, `words`, `n` |
| `POST /api/z04/top20` | `text`, `n` | 200 `source`, `words`, `n` |
| `POST /api/z06/plot` | `expr`, `xmin`, `xmax`, `samples` | 200 `segments`, `asymptotes`, `yclamp`, `samples`, `samplesRequested` |
| `POST /api/z07/solve` | `equation`, `limit` | 200 `count`, `truncated`, `solutions`, `spaceSize`, `timeMs` |
| `POST /api/z08/generate` | `seed` (0 → current time) | 200 `puzzle`, `solution`, `blanks`, `removed`, `timeMs` |
| `POST /api/z08/solve` | `grid` (81 values) | 200 `solved`, `solution`, `count`, `capped`, `timeMs` |

# Rules the port guarantees

- Every input bound is checked before work starts; failure is a 400 with a stated range
  ([Synchronous validation](/rules/synchronous-validation.md)).
- Applied defaults and applied clamps are echoed
  ([Explicit defaults](/rules/explicit-defaults.md),
  [No silent truncation](/rules/no-silent-truncation.md)).
- Requests are body-size limited (4 MB).
- A capability that finds nothing answers 200 with an empty/zero result; "no solution" is data.
- `POST /api/z08/generate` failure is an internal error, not a client error — the client sent
  nothing that can be wrong beyond the seed, which is unbounded by design.

# Timeouts

The server's write timeout is generous (2 minutes) precisely because the sudoku solve and the
cryptarithm search can be long; a synchronous port is only used where the worst case is bounded.

# Citations

- `cmd/wdi-server/main.go` — the endpoint handlers and their validation.
- `cmd/wdi-server/api_test.go`, `handlers_test.go`, `z08_ui_test.go`.
- `internal/z01`, `z04`, `z06`, `z07`, `z08` — the answering capabilities.
