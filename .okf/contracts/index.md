# contracts

* [Background job API](/contracts/background-job-api.md) - The shared submit/poll/cancel contract for the three long-running capabilities, with the jobId, progress snapshot and cancellation shapes.
* [Error response contract](/contracts/error-response.md) - The single failure envelope for the API: a JSON error object, a cause-appropriate status code, and Polish messages naming the value and the range.
* [Solver capability port](/contracts/solver-package-port.md) - The internal port between serving layer and solvers: pure and cancellable-with-progress signatures, replaceable call points, and the obligations every capability honours.
* [Static asset delivery](/contracts/static-assets.md) - Delivery of the embedded single-page client (markup, stylesheet, script) as static content alongside the API.
* [Synchronous solver API](/contracts/synchronous-solver-api.md) - The synchronous JSON port for the bounded-cost capabilities: request shapes, success payloads, validation and echo of applied values.
* [Task catalog API](/contracts/task-catalog-api.md) - Read-only contract exposing the eight task descriptors, consumed by the client, the container healthcheck and the gate probe.
