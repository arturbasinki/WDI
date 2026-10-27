---
type: QualityAttribute
title: Accessibility (WCAG 2.1 AA)
description: "Usability/accessibility at WCAG 2.1 AA: measured text contrast of at least 4.5:1 (5.44-6.39:1 achieved), every label bound, every button typed, announced progress."
tags: [quality-attribute, iso25010, usability, accessibility]
timestamp: "2026-10-27T16:32:44Z"
---

# Attribute class (ISO/IEC 25010)

**Usability** — accessibility, specifically WCAG 2.1 level AA for the text the interface renders.

# Criteria in force

| Criterion | Requirement | Achieved |
|---|---|---|
| Text contrast | ≥ 4.5:1 for body and auxiliary text | 18 token pairs at 5.44–6.39:1, the page-background pair at 6.39:1 |
| Label association | every form label bound to its field | 18 of 18 labels carry the binding |
| Button semantics | explicit button type in forms | 37 of 37 buttons typed |
| Progress indication | announced progress with a current value | progress indicators expose role and current value |
| Language | all copy in Polish, document language declared | yes |

# How it is kept

The contrast pairs of the stylesheet are computed by a script and the check exits non-zero on a
failure, so it is a gate step and not an opinion. The auxiliary text token was measured at
4.17:1 (a failure) and darkened to reach the range above.

# Design constraints that follow

Text colours are tokens, not literals, so a contrast fix is one change. Progress bars must be
announced, not merely drawn. Chart bars encode values with a log scale and a text label, so the
information is not colour-only.

# Related

[Browser presentation](/domains/browser-presentation.md) ·
[Browser client context](/contexts/browser-client.md) ·
[Error communication](/rules/error-communication.md)

# Citations

- `.verify/contrast.py` — the token-pair contrast computation, exit code as the gate.
- `cmd/wdi-server/web/styles.css` — the token definitions incl. the corrected muted text token.
- `cmd/wdi-server/web/app.js`, `index.html` — labels, typed buttons, progress attributes,
  declared language.
- `README.md` — the accessibility section stating the achieved ratios.
