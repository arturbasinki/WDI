---
type: Domain
title: Background computation (domain)
description: "The supporting domain that turns long-running computations into identified, pollable, cancellable, resource-bounded background work."
tags: [domain, supporting, jobs]
timestamp: "2026-10-27T16:32:44Z"
---

# Kind of subdomain

Supporting, and the heart of the serving layer. Three capabilities take minutes, not
milliseconds, so they are executed as observable background work rather than inside a request.

# Responsibility

Give a long computation an identity, a status lifecycle, a progress channel, a cancellation
handle, and a bounded, self-cleaning registry — while guaranteeing that the machine is never
oversubscribed.

# Concepts it owns

- [Job](/model/job.md) — the run, its status, progress and result.
- [Job lifecycle workflow](/workflows/job-lifecycle.md) — accepted → running → terminal.
- [Background job policy](/rules/background-job-policy.md) — validation, concurrency cap,
  cancellation, retention.
- Cancellation is cooperative: the computation observes a cancellation signal at a chosen
  granularity and returns the partial state together with the cancellation
  ([Cancellation granularity](/rules/cancellation-granularity.md)).

# Why it is a domain and not plumbing

Its rules are business rules, not incidental mechanics: a finished job must not report
"cancelling", an over-submitted request is refused rather than queued, expired jobs vanish
and become 404, and a cancellation must not hang. Each of these is asserted by a test that
pins the HTTP-observable behaviour.

# Related

[HTTP API context](/contexts/http-api.md) ·
[Background job API](/contracts/background-job-api.md) ·
[Detached job context decision](/decisions/detached-job-context.md) ·
[Lazy retention decision](/decisions/lazy-retention-no-timers.md)

# Citations

- `cmd/wdi-server/main.go` — the job record (mutex-guarded fields), the registry with lazy
  sweep and channel semaphore, the job start helper (incl. panic recovery), the shared
  GET/DELETE handlers, and the shutdown-time cancellation of all live jobs.
- `cmd/wdi-server/handlers_test.go` — 400-before-job-creation, 429 at the third concurrent
  submission, 202→cancelled→409, TTL expiry with an injected clock, capacity ceiling.
