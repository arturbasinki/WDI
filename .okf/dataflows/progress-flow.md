---
type: DataFlow
title: Progress reporting flow
description: "How progress travels from a computation loop through a lock-protected job snapshot to a polling client, with copies, a single latest snapshot, and defined poll stops."
tags: [dataflow, jobs, progress]
timestamp: "2026-10-27T16:32:44Z"
---

# Path

Computation loop → job snapshot (locked) → poll → progress bar / partial results in the view.

# Stages

1. **Emit** — the capability calls its progress callback at its cadence: ≈10 Hz for the number
   search, per 500 000 keys for the cipher attack, per completed measurement for the benchmark
   ([Progress cadence](/rules/progress-cadence.md)).
2. **Publish** — the payload is a **copy** (sorted where order matters) handed to the job,
   stored under the job's lock as the latest snapshot; the producer is joined before the final
   result is assembled.
3. **Store** — the job keeps only the latest snapshot, not a history. Nothing is written to disk.
4. **Poll** — the client fetches the snapshot on an interval; the response carries
   `jobId`, `status`, and `progress` when present.
5. **Render** — numbers found so far, processed/limit ratio, keys tried, or measurements
   completed.
6. **Terminate** — the poller stops on a terminal status, on a fetch failure (with a visible
   banner), and on leaving the view.

# Failure modes this flow prevents

- A progress producer racing the final assembly (caught by the race detector).
- A `done: true` emitted mid-attack, which stops the client polling on an incomplete answer.
- A poller swallowing a fetch failure, leaving a progress bar running forever.

# Related

[Search progress](/model/search-progress.md) · [Job](/model/job.md) ·
[Background job API](/contracts/background-job-api.md) ·
[Progress cadence](/rules/progress-cadence.md)

# Citations

- `internal/z02/armstrong.go`, `internal/z03/sieve.go`, `internal/z05/gronsfeld.go` — emitters.
- `cmd/wdi-server/main.go` — `setProgress`, `snapshot`.
- `cmd/wdi-server/web/app.js` — `poll`, `startJob`, `jobBox`.
