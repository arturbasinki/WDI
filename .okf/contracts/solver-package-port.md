---
type: Contract
title: Solver capability port
description: "The internal port between serving layer and solvers: pure and cancellable-with-progress signatures, replaceable call points, and the obligations every capability honours."
tags: [contract, port, hexagonal, testing]
timestamp: "2026-10-27T16:32:44Z"
---

# Port

The internal seam between the serving layer and a solver capability. It is a port in the
hexagonal sense: the serving layer depends on the capability's signature, never on its
internals, and the heavy three are invoked through replaceable call points so HTTP behaviour
can be tested without real computation.

# Signatures the port admits

| Kind | Shape | Used by |
|---|---|---|
| Pure computation | `inputs → (result, error)` | z01, z04, z06, z07, z08 |
| Cancellable, observed search | `(cancellation, params, progressCallback) → (result, error)` | z02, z03, z05 |

# Obligations of any capability behind this port

- A nil progress callback is accepted and means "no progress".
- Cancellation returns the **partial** state together with the cancellation; never a success.
- Invalid parameters are errors, never substituted defaults
  ([Explicit defaults](/rules/explicit-defaults.md)).
- Only mathematical bounds live here; user-facing bounds belong to the web layer
  ([Layer ownership of bounds](/rules/layer-ownership-of-bounds.md)).
- No panics escape: the caller converts one into a job error.
- Returned collections are copies, not aliases of internal state
  ([Search progress](/model/search-progress.md)).
- Reference content (canonical text, dictionary) is carried by the capability
  ([Embedded corpus](/model/embedded-corpus.md)).

# Replaceability

The serving layer calls the three heavy capabilities through named indirections; tests swap
them for stubs, and the test harness drains running job goroutines before restoring the real
ones, so a stale goroutine cannot read a stub after its test ended.

# Citations

- `internal/z02/armstrong.go`, `internal/z03/sieve.go`, `internal/z05/gronsfeld.go` — the
  `(ctx, params, onTick)` signatures with nil guards and cancel contracts.
- `cmd/wdi-server/main.go` — the named indirections (`searchArmstrong`, `runBenchmark`,
  `breakGronsfeld`).
- `cmd/wdi-server/handlers_test.go` — seam substitution and the drain-before-restore ordering.
