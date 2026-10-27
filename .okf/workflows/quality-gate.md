---
type: Workflow
title: Quality gate workflow
description: "The test-first verification sequence: format, vet, race-enabled tests, frontend syntax, build, live probes, contrast check — ending in a named green state."
tags: [workflow, quality, process]
timestamp: "2026-10-27T16:32:44Z"
---

# Scope

The verification process a change must pass before it is accepted, and the discipline around it.

# Steps

1. **Reproduce as a test first.** A reported defect becomes a red regression test that pins the
   HTTP-observable or numeric behaviour, one test per acceptance criterion.
2. **Fix**, keeping the layer boundaries intact
   ([Layer ownership of bounds](/rules/layer-ownership-of-bounds.md)).
3. **Run the fast checks**: source formatting (scoped to the source trees), static analysis.
4. **Run the full suite** with the data-race detector, over the solvers **and** the serving
   layer; the suite is required to be fast (seconds in short mode, tens of seconds in full),
   which is enforced by budgeting the heavy cases
   ([Test cost budget decision](/decisions/test-cost-budget.md)).
5. **Check the frontend's syntax** — it is one unbundled script, so syntax is the checkable
   surface.
6. **Build the artifact** with the same flags as the image, and **build the image**.
7. **Probe the running server live**: the catalog returns eight tasks; an out-of-range search
   limit is a 400; the largest linear system answers within a size budget; a million nested
   parentheses is a 400 with the server still alive; SIGTERM exits 0.
8. **Check accessibility numerically**: the token pairs of the stylesheet must meet 4.5:1.
9. Green across all of it is the accepted state; the outcome is named, not implied.

# Discipline around the steps

- Multi-agent repair work is sequenced (fix → review → next), never parallel: the work touches
  the same few files and shares one working tree.
- Commits happen only on request.
- Numeric claims in the documentation are verified through the live API, never by reading code
  ([Reproducibility](/quality/reproducibility.md)).

# Related

[Maintainer](/actors/maintainer.md) · [Hosted CI runner](/external-systems/ci-runner.md) ·
[Workspace cache constraint](/constraints/workspace-cache.md)

# Citations

- `.verify/gate.sh` — the independent gate incl. the live probes, ending `BRAMKA_OK`.
- `.verify/contrast.py`, `.verify/p1-fe-check.sh`, `.verify/shot.sh`, `.verify/cdp.mjs`.
- `.github/workflows/ci.yml` — the same bar as pipeline steps, Go version from the module file.
- `RAPORT-JAKOSC-KODU.md` §9 — the environment briefing and the sequenced repair plan.
