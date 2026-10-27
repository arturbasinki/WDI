---
type: Workflow
title: Graceful shutdown workflow
description: "On a stop signal the server cancels every live job, closes the listener within 5 seconds and exits with status 0."
tags: [workflow, shutdown, deployment]
timestamp: "2026-10-27T16:32:44Z"
---

# Scope

What happens when the platform stops the process, and what the platform is entitled to expect.

# Steps

1. The platform (or a terminal) delivers a stop signal — termination or interrupt.
2. Every live job is cancelled. This is the point: computations that outlived their client are
   stopped deliberately, not abandoned.
3. The listener is closed with a grace window of at most 5 seconds for in-flight requests.
4. The process waits for the shutdown goroutine to complete, then exits with status 0.

# Why it matters here

The container platform's stop has its own grace period; a process that ignores the signal is
killed at the platform's deadline, and its logs show a non-zero exit. Honouring the signal
means the container stops immediately and cleanly. Job state is ephemeral by design
([Job retention](/rules/job-retention.md)), so nothing needs to be persisted — only
interrupted.

# Observable acceptance

A live probe sends the termination signal and asserts exit status 0; it is part of the
[Quality gate workflow](/workflows/quality-gate.md).

# Related

[Container platform](/external-systems/container-platform.md) ·
[Job](/model/job.md) · [Job lifecycle workflow](/workflows/job-lifecycle.md)

# Citations

- `cmd/wdi-server/main.go` — the signal-notify context, cancelling all jobs, the 5-second
  shutdown timeout, waiting for the shutdown goroutine, exit 0.
- `.verify/gate.sh` — the SIGTERM probe asserting exit 0.
- `docker-compose.yml` — restart policy that relies on clean exits.
