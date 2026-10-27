# quality

* [Accessibility (WCAG 2.1 AA)](/quality/accessibility.md) - Usability/accessibility at WCAG 2.1 AA: measured text contrast of at least 4.5:1 (5.44-6.39:1 achieved), every label bound, every button typed, announced progress.
* [Portability and single-artifact delivery](/quality/portability.md) - Portability (ISO 25010): one self-contained artifact, no configuration, no data files, no external services, buildable offline; the bundle describes it language-agnostically.
* [Race freedom](/quality/race-freedom.md) - Race freedom: the full suite runs under a data-race detector in the gate; job state, progress producers, registry, slots and test seams are lock- and join-disciplined.
* [Reliability and crash containment](/quality/reliability.md) - Reliability: no documented request can crash the process or exhaust memory — bounded inputs, narration, recursion, concurrency, crash containment, overflow guards, and live probes.
* [Reproducibility of results](/quality/reproducibility.md) - Reproducibility: every documented numeric result is pinned by explicit-value tests and re-verified against a live process, with determinism engineered into seeded and ranked output.
* [Resource boundedness](/quality/resource-boundedness.md) - Performance efficiency: every resource has an explicit bound (input sizes, narration cells, jobs, slots, timeouts), and the reach claims are measured, not assumed.
* [Security posture](/quality/security-posture.md) - Security posture: no XSS, no command execution, no outbound calls, no traversal, bounded resource abuse, non-root and loopback-only deployment, with accepted gaps recorded.
* [Test determinism](/quality/test-determinism.md) - Test determinism: explicit full expected values, seeded and ordered output, injected clocks, stubbed heavy computations, and job goroutines drained between tests.
