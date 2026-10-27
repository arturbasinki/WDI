---
type: Constraint
title: Standard library only
description: "The product uses only the language's standard library and the browser's built-ins: no third-party dependencies, no frontend build step."
tags: [constraint, arc42-2, dependencies]
timestamp: "2026-10-27T16:32:44Z"
---

# Constraint

No third-party dependencies. The product is built from the language's standard library and the
browser's built-in capabilities; the client uses no framework, no bundler and no build step.

# Type and origin

Technical, self-imposed, and reaffirmed after the rewrite from Python. It is enforced
mechanically: the module listing shows a single module, and the CI gate builds the binary and
the image from that state.

# What it forces

- The expression language, the chart drawing, the progress machinery, the counting and ranking
  utilities are first-party code
  ([Single-binary decision](/decisions/stdlib-only-single-binary.md)).
- The client's HTML/JS/CSS must be checkable as plain text — the gate runs a syntax check on
  the single script.
- The runtime environment needs no package installation, which is what lets the container's
  runtime stage build offline
  ([Two-stage image decision](/decisions/multistage-image.md)).
- A deeper static analyser is excluded from the gate because it would have to be downloaded and
  the environment cannot fetch it.

# What must not be done because of it

Do not add a framework to the client, a plotting or algebra dependency to the solver, a database
for job state, or a task runner that needs installation. Each of those changes the deployment
story, which is the point of the constraint.
