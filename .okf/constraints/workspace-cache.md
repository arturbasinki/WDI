---
type: Constraint
title: Workspace-local build cache
description: "The working environment keeps build cache, toolchain workspace and container-CLI home inside the repository (read-only home), forcing scoped checks and local-only verification scripts."
tags: [constraint, environment, tooling]
timestamp: "2026-10-27T16:32:44Z"
---

# Constraint

The working environment's home directory is read-only, so the build cache, the module/tool
workspace, and a writable home for the container CLI must live **inside** the repository
directory, exported on every command.

# Consequences for tooling and for the bundle

- Every command that touches the toolchain runs with the cache and workspace set to repository
  subdirectories; this is roughly 230 MB of local artefacts.
- Formatting and static-analysis checks must be **scoped to the source trees**. A repository-wide
  scan walks the toolchain and cache directories — part of which is deliberately not parseable —
  and fails with noise instead of findings.
- These artefacts are build environment, not source, and are git-ignored; the container build
  context excludes them.
- The local verification scripts hard-code these paths, so they are environment-specific tooling
  rather than a portable part of the product.

# What this is not

Not a property of the product. Deployed elsewhere, the cache lives in the normal location and
the checks can be repository-wide
([Hosted CI runner](/external-systems/ci-runner.md)).

# Citations

- `.verify/gate.sh` — the exported cache and workspace paths, the scoped formatting check.
- `.gitignore` — the cache, toolchain, binary and screenshot directories.
- `.dockerignore` — the same directories excluded from the build context.
- `RAPORT-JAKOSC-KODU.md` §5 P3 item 20 and §9.1 — the artefact size and the environment
  briefing.
