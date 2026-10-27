# contexts

* [Browser client context](/contexts/browser-client.md) - The bounded context of the delivered single-page client: routing, views, polling, rendering, and disciplined release of interactive resources.
* [HTTP API context](/contexts/http-api.md) - The bounded context serving the product: transport, edge validation, job policy, task catalog, static delivery, and test seams over the heavy computations.
* [Solver library context](/contexts/solver-library.md) - The bounded context of pure computational solvers: one module per task, no transport, no shared state, errors for invalid input, cooperative cancellation.
