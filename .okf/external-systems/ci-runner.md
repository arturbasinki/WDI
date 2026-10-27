---
type: ExternalSystem
title: Hosted CI runner
description: "The hosted pipeline environment that enforces the project's declared quality bar on every push and pull request."
tags: [external, ci]
timestamp: "2026-10-27T16:32:44Z"
---

# Role

Runs the project's quality gate automatically. The steps are ordered cheap → expensive and
exist to keep the claims made in `README.md` from drifting: formatting, static analysis,
full test suite with the data-race detector, frontend syntax check, binary build, image build.

# Boundaries the runner must respect

- The formatting check is deliberately scoped to the source trees only. A repo-wide scan
  walks local toolchain and build-cache directories, part of which is intentionally not
  parseable, and fails the gate with noise — see [Workspace cache constraint](/constraints/workspace-cache.md).
- The Go version is taken from the module definition, so the gate cannot drift from the
  build image.
- A deeper static analyser is intentionally absent: it must be downloaded and the working
  environment has no access to a module proxy.

# Related

[Quality gate workflow](/workflows/quality-gate.md) ·
[Test cost budgeting decision](/decisions/test-cost-budget.md)

# Citations

- `.github/workflows/ci.yml` — six steps, `go-version-file`, scoped `gofmt`, `go vet`,
  `go test -race -count-1`, `node --check`, binary build, `docker build`.
