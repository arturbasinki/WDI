---
type: BusinessRule
title: No silent truncation
description: "Responses disclose every silent adjustment: requested vs applied counts, applied caps, dropped characters, truncation flags — never a plausible-looking lie."
tags: [business-rule, honesty, api]
timestamp: "2026-10-27T16:32:44Z"
---

# Rule

Any parameter the system applies differently from what was requested, and any input content it
discards, must be reported in the response. A response states what was **used**, alongside what
was **asked for**.

# Instances of the rule

| Capability | What is disclosed |
|---|---|
| z06 plot | `samplesRequested` next to `samples` (clamped to 501…20001); the applied `yclamp` |
| z06 errors | unknown characters named as whole characters, not raw bytes |
| z05 cipher | `skipped` — the unique characters dropped by normalization (omitted when none) |
| z05 attack | `capApplied` — the key cap was clamped, so the search stopped at the cap, not at the end of the key space |
| z07 | full `count` next to a `solutions` list cut to `limit`, plus `truncated` |
| z04 | the `n` actually used, in both request shapes |
| z01 | `stepsTruncated` — the narration was cut by the memory budget |

# Why

Every one of these was found as a silent behaviour in an audit: a plot with branches ending
mid-air, a cipher answer missing seven letters, an attack "finishing" at a cap, a "top 50" list
presented as if the puzzle had 50 solutions, a request for 5 000 000 samples answered with a
plausible-looking 1601. Silence reads as correctness, so it is a defect, not a cosmetic gap.

# Consequence for new capabilities

Adding a capability means deciding, up front, what it may silently adjust — and answering it
with a field. A rejection with a stated range is acceptable; a silent adjustment is not.

# Related

[No-silent-adjustment for bounds](/rules/layer-ownership-of-bounds.md) ·
[Solution capping](/rules/solution-capping.md) ·
[Normalization report](/model/normalization-report.md) ·
[Plot clamping](/model/plot-clamping.md) ·
[Reproducibility of results](/quality/reproducibility.md)

# Citations

- `cmd/wdi-server/main.go` — `skipped` emitted only when non-empty; `n` echoed.
- `internal/z05/gronsfeld.go` (`capApplied`, `Skipped`), `internal/z06/expr.go`
  (`SamplesRequested`, `YClamp`, rune-correct message), `internal/z07/cryptarithm.go`
  (`Count`, `Truncated`), `internal/z04/words.go`, `internal/z01/solver.go` (`StepsTruncated`).
- `cmd/wdi-server/web/app.js` — the client surfaces the dropped-characters warning.
