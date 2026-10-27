---
type: BusinessRule
title: Error communication
description: "Failures are JSON error objects with an appropriate status, phrased in Polish and naming the offending value and the accepted range."
tags: [business-rule, api, i18n, convention]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

Every failure is communicated as a single-kind JSON object `{ "error": "<message>" }` with an
appropriate status code, and the message is written **in Polish**, naming the offending value
and the accepted range.

# Consequences

- A learner reading the interface understands the rejection without consulting documentation.
- A message that states the range makes the boundary testable from the UI: "limit must be in
  1000…1000000000000 (10¹²)".
- A parser message identifies the position and the token: row number for a matrix coefficient,
  token for an expression, cell index (1-based) for a sudoku grid, character for an unknown
  symbol.
- A rejection that lists what *is* available is preferred over one that only says "unknown"
  (the expression parser names the available variables and functions).
- Messages are counted in characters (runes), so a non-ASCII character is reported intact; a
  byte-wise report rendered "π" as mojibake.

# Boundary of the rule

Polling a job whose computation failed is a 200 with the failure inside the job snapshot; the
job answered. The failure of the *job* is data, not a transport error
([Background job policy](/rules/background-job-policy.md)).

# Related

[Error response contract](/contracts/error-response.md) ·
[Synchronous validation](/rules/synchronous-validation.md) ·
[Browser presentation](/domains/browser-presentation.md)

# Citations

- `cmd/wdi-server/main.go` — the JSON success/error helpers used by every handler.
- `internal/z01/solver.go`, `internal/z06/expr.go` (rune-correct `utf8.DecodeRuneInString`
  message, `funcNames()`), `internal/z07/cryptarithm.go` (token-count and operator messages),
  `internal/z08/sudoku.go` (cell-indexed validation messages).
- `cmd/wdi-server/handlers_test.go` — asserts on message text (e.g. the Polish conflict).
