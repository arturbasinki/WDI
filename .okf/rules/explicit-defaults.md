---
type: BusinessRule
title: Explicit defaults
description: Defaults are applied by the calling layer and echoed in responses; a library rejects a non-positive count rather than silently substituting a default.
tags: [business-rule, api, convention]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

A default value is applied by the **caller**, explicitly. A library never substitutes a default
for a value the caller supplied — in particular, a zero or negative "how many" is an error, not
an invitation to fill in the default.

# Instances

| Default | Value | Applied by | Library behaviour on ≤ 0 |
|---|---|---|---|
| top words (`DefaultTop`) | 20 | the handler | error: "k must be ≥ 1" |
| narcissistic limit (`DefaultLimit`) | 10¹⁰ | the handler (when the request omits it) | accepted: 0 is out of range for the web layer, not for the maths |
| cryptarithm solutions (`DefaultSolutions`) | 50 | the handler | means "no limit" |
| cipher key cap (`BreakCapDefault`) | 2 000 000 | the handler (when omitted) | means "use the default" |
| plot samples | 1601 | the handler | clamped to 501…20001 and disclosed |

# Why

Substituting a default inside the library makes the effective value invisible and unaskable:
the caller cannot request "zero results", cannot tell which value was used after the fact, and
the API response can present a default as if it were the request. The response now echoes the
`n` that was used.

# Related

[Layer ownership of bounds](/rules/layer-ownership-of-bounds.md) ·
[No silent truncation](/rules/no-silent-truncation.md) ·
[Solution capping](/rules/solution-capping.md)

# Citations

- `internal/z04/words.go` — `DefaultTop` documented as caller-applied; `TopWords` rejects
  `k < 1`.
- `internal/z07/cryptarithm.go` — `DefaultSolutions`/`MaxSolutions`, `limit <= 0` = unlimited.
- `internal/z05/gronsfeld.go` — `BreakCapDefault`, `BreakCapMax`, `CapApplied`.
- `cmd/wdi-server/main.go` — defaults applied in handlers, `n` echoed in the response.
