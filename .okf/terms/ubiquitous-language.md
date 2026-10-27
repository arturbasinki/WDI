---
type: Term
title: Ubiquitous language
description: "The project's vocabulary: task, job, task catalog, seam, gate, honest response, and the domain words used across solver, API and frontend."
tags: [glossary, ubiquitous-language]
timestamp: "2026-10-27T16:32:44Z"
---

# Vocabulary

**Task (zadanie, z01…z08)** — one exercise of the course set; the unit of both the domain
model and the user-facing catalog. Each task is an independent computational module with
its own tests and its own view.

**Task catalog** — the single registry of task identities, numbers, titles and blurb
([Task descriptor](/model/task-descriptor.md)). The frontend keeps no copy.

**Job** — one background execution of a long-running computation, identified by a job id,
with a status lifecycle and a pollable progress record
([Job aggregate](/model/job.md)). Not a task; a task is a capability, a job is a run.

**Seam** — a swappable call point in the serving layer through which the heavy computations
are invoked, so HTTP tests can simulate long jobs without real computation.

**Gate (bramka)** — the verification pass that must end green before a change is accepted;
its positive outcome is named `BRAMKA_OK`.

**Honest response ("prawda API")** — a response that reports what the server actually used:
requested vs. applied sample counts, applied key caps, skipped characters, full solution
counts alongside truncated lists
([No silent truncation](/rules/no-silent-truncation.md)).

**Verified result** — a numeric claim in the documentation that has been reproduced by
running the system.

## Domain words

| Term | Meaning | Where |
|---|---|---|
| Leading element (pivot) | Largest-magnitude row chosen for elimination | [Linear systems](/domains/linear-systems.md) |
| RREF | Reduced row echelon form, the elimination target | [Linear systems](/domains/linear-systems.md) |
| Contradictory / underdetermined / unique | The three classifications of a linear system | [Solution status](/model/solution-status.md) |
| Digit-power sum | Sum of digits raised to the digit-count power | [Narcissistic numbers](/domains/narcissistic-numbers.md) |
| Chunk / block | Ten consecutive numbers sharing the leading part | [Narcissistic numbers](/domains/narcissistic-numbers.md) |
| Structure variant | Slice / map / bitmap implementation of the sieve | [Sieve benchmark row](/model/sieve-benchmark-row.md) |
| Cipher alphabet | Space + 26 lowercase letters, 27 symbols | [Cipher key](/model/cipher-key.md) |
| Normalization | Lower-casing plus dropping symbols outside the alphabet | [Normalization report](/model/normalization-report.md) |
| Dictionary hit ratio | Fraction of decrypted words present in the dictionary (> 80% ⇒ key found) | [Cipher break search](/workflows/cipher-break-search.md) |
| Implied multiplication | `2x` means `2*x` | [Function plotting](/domains/function-plotting.md) |
| Outlier sample | Sample whose \|y\| exceeds the IQR-derived threshold | [Plot clamping](/model/plot-clamping.md) |
| Asymptote | Detected sign flip with a large jump; ends a segment | [Plot segment](/model/plot-segment.md) |
| Letter / no-zero letter | Symbol of a cryptarithm / letter forbidden to be 0 | [Cryptarithm puzzle](/model/cryptarithm-puzzle.md) |
| Digging | Removing cells from a full grid while uniqueness holds | [Sudoku puzzle](/model/sudoku-puzzle.md) |
| MRV heuristic | Next cell = the one with fewest candidates | [Sudoku solver workflow](/workflows/sudoku-solve.md) |
| Loopback binding | Publishing the port on the local interface only | [Container platform](/external-systems/container-platform.md) |

# Citations

- `cmd/wdi-server/main.go` — `Job`, `jobStore`, `taskCatalog`, test seams (`searchArmstrong`, `runBenchmark`, `breakGronsfeld`).
- `.verify/gate.sh` — `BRAMKA_OK` / `BRAMKA_FAIL`.
- `internal/z01..z08` — package doc comments carry the domain vocabulary in Polish.
