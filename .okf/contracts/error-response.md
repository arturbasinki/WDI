---
type: Contract
title: Error response contract
description: "The single failure envelope for the API: a JSON error object, a cause-appropriate status code, and Polish messages naming the value and the range."
tags: [contract, api, errors]
timestamp: "2026-10-27T16:32:44Z"
---

# Port

Every failure of the API has one shape and one language.

# Shape

```json
{ "error": "<Polish sentence naming the value and the accepted range>" }
```

with the status code chosen by cause:

| Cause | Status |
|---|---|
| Parameter outside its stated range, malformed input, unparseable body | 400 |
| Unknown job id (never existed, or expired from retention) | 404 |
| Cancellation requested on an already finished job | 409 |
| Both computation slots busy | 429 |
| Internal generation failure | 500 |

# Rules

- Messages are Polish, since the interface is Polish and the learner reads them
  ([Error communication](/rules/error-communication.md)).
- A range rejection states the range and the offending value, so the boundary is discoverable
  from the UI.
- A parse rejection names a position: row, cell index, token, character, function arity,
  nesting depth.
- A character is named as a character, not as a byte, so non-ASCII input reads correctly.
- An unknown name lists what is available.
- An over-long body is rejected with 400 (recorded as a deliberate deviation from the
  semantically correct 413).

# Non-goals

A job whose computation failed is **not** an error response: polling it is 200 with
`status: "error"` inside the snapshot ([Background job policy](/rules/background-job-policy.md)).

# Citations

- `cmd/wdi-server/main.go` — the JSON error helper used by every handler; the 4 MB body limit.
- `cmd/wdi-server/handlers_test.go` — assertions on status **and** message text.
- `RAPORT-JAKOSC-KODU.md` §5 P3 — the 400-vs-413 deviation.
