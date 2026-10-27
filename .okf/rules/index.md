# rules

* [Animation step budget](/rules/animation-step-budget.md) - Recorded elimination steps are bounded by 400,000 matrix cells, with a truncation flag, so animation narration cannot exhaust memory.
* [Background job policy](/rules/background-job-policy.md) - The complete status-code policy for background jobs: 202 accept, 400 validation, 429 capacity, 404 unknown, 202 cancelling, 409 finished, 200 snapshots.
* [Cancellation granularity](/rules/cancellation-granularity.md) - Cancellation is observed inside long loops at a near-zero-cost granularity (once per 2^20 iterations) so a slow measurement is interrupted in a fraction of a second.
* [Cancellation semantics](/rules/cancellation-semantics.md) - Cancellation semantics by job state: running cancels with an accepted acknowledgement, finished conflicts, unknown is not found; the client also cancels on leaving a view.
* [Computation concurrency cap](/rules/concurrency-cap.md) - At most two computations run concurrently; the third submission is refused with a too-many-requests response and no job is created.
* [Error communication](/rules/error-communication.md) - Failures are JSON error objects with an appropriate status, phrased in Polish and naming the offending value and the accepted range.
* [Explicit defaults](/rules/explicit-defaults.md) - Defaults are applied by the calling layer and echoed in responses; a library rejects a non-positive count rather than silently substituting a default.
* [Job retention policy](/rules/job-retention.md) - Finished jobs are kept 15 minutes and at most 200 at a time, evicted lazily on access, oldest finished first, never a running job.
* [Layer ownership of bounds](/rules/layer-ownership-of-bounds.md) - User-facing bounds live in the web layer; the solver library enforces only mathematical limits, so small and large inputs are both honest there.
* [No silent truncation](/rules/no-silent-truncation.md) - Responses disclose every silent adjustment: requested vs applied counts, applied caps, dropped characters, truncation flags — never a plausible-looking lie.
* [Progress cadence](/rules/progress-cadence.md) - Progress is published at a fixed, low-cost cadence (10 Hz, per 500k keys, per measurement) with copies, a done flag reserved for finished runs, and joined producers.
* [Solution capping](/rules/solution-capping.md) - A solution limit trims the returned list only; the full solution count and the truncation flag are always reported.
* [Synchronous validation](/rules/synchronous-validation.md) - Out-of-bounds parameters are rejected at submission time with a client error, before any job exists and before a computation slot is taken.
* [View resource disposal](/rules/view-resource-disposal.md) - Views register every timer, listener and poller they start and release them on route change; a stale route may not overwrite a newer one.
* [XSS escaping](/rules/xss-escaping.md) - All externally sourced values are HTML-escaped before reaching the DOM; no raw markup is built from foreign data.
