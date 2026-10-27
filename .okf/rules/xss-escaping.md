---
type: BusinessRule
title: XSS escaping
description: All externally sourced values are HTML-escaped before reaching the DOM; no raw markup is built from foreign data.
tags: [business-rule, security, frontend]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

Every value that reaches the page from outside the client — API results, error messages, user
input echoed back, embedded corpus content — is inserted as escaped text. Raw markup assembly
from foreign data is prohibited.

# Scope verified

The verification was exhaustive over all eight views: word-frequency rows, the echoed plot
formula in the form, cipher text and key, cryptarithm equations, API error strings, progress
messages. No interpolated markup path from source data to the DOM exists.

# Why it is a standing rule

The client builds views from template strings, so an escaping helper is one call away from being
forgotten, and the data is exactly the kind a learner types: equations, formulas, ciphers. The
rule is cheap, absolute and testable by inspection; the audit records it as "XSS: none".

# Related safety rules

The server makes no outgoing requests, embeds no commands, and cannot be traversed outside its
embedded content; the request body is size-limited; process crashes inside a computation become
a job error, not a dead server.

# Related

[Browser client context](/contexts/browser-client.md) ·
[Security posture](/quality/security-posture.md) ·
[Error response contract](/contracts/error-response.md)

# Citations

- `cmd/wdi-server/web/app.js` — the `esc` helper applied at every insertion point.
- `RAPORT-JAKOSC-KODU.md` §4 — the exhaustive XSS verification and the surrounding security notes.
- `cmd/wdi-server/main.go` — the 4 MB body limit; panic recovery in the job worker.
