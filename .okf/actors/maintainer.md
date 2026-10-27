---
type: Actor
title: Maintainer / auditor
description: The person (or agent) who changes the code and is accountable for the documented claims staying true; runs the quality gate.
tags: [actor, quality]
timestamp: "2026-10-27T16:32:44Z"
---

# Role

Evolves the solvers and the serving layer, audits them for correctness and resource
safety, and is bound by the project's non-negotiable constraints
([Standard-library-only constraint](/constraints/stdlib-only.md)). Changes are made
test-first: a red regression test for the reported defect, then the fix, then the gate.

# Interactions

- Runs the [Quality gate workflow](/workflows/quality-gate.md) before a change is accepted.
- Writes regression tests that pin the job-policy contract
  ([Job policy rules](/rules/background-job-policy.md)) one-to-one with acceptance criteria.
- Verifies numeric claims through the live API, not by reading code.
- Commits only on request.

# Citations

- `.verify/gate.sh` — the reviewer's independent gate (`BRAMKA_OK` / `BRAMKA_FAIL`).
- `.verify/contrast.py` — accessibility contrast checker.
- `RAPORT-JAKOSC-KODU.md` — audit report: findings, fixes, open items, environment constraints.
- `cmd/wdi-server/handlers_test.go` — HTTP-layer regression tests, incl. injected clock.
