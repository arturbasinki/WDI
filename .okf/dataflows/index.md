# dataflows

* [Browser client flow](/dataflows/browser-client-flow.md) - How the client moves between views: route change, disposal of the previous view's resources, catalog-driven naming, escaped rendering, polling, and the route race guard.
* [Corpus to ranking flow](/dataflows/corpus-to-ranking-flow.md) - How embedded reference content reaches its consumers: the literary corpus through tokenising, filtering, tallying and ranking, and the dictionary as the cipher attack's oracle.
* [Benchmark measurement flow](/dataflows/measurement-flow.md) - How benchmark data is produced and consumed: fixed structure order per range, collection-forced timings, rows, progress ticks, and log-scaled rendering.
* [Progress reporting flow](/dataflows/progress-flow.md) - How progress travels from a computation loop through a lock-protected job snapshot to a polling client, with copies, a single latest snapshot, and defined poll stops.
* [Request validation flow](/dataflows/request-validation-flow.md) - How a request travels from the browser form through intake, decoding, edge validation and interpretation to a capability, and where each rejection branch exits.
* [System composition flow](/dataflows/system-composition.md) - How the artifact is assembled from embedded content and how data moves between browser, serving layer, solvers, jobs and the container platform at run time.
