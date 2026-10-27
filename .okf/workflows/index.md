# workflows

* [Cipher break search workflow](/workflows/cipher-break-search.md) - Enumerate digit keys by increasing length, accept the decryption whose words are over 80% dictionary-covered, with early abandon, key caps and cancellation.
* [Graceful shutdown workflow](/workflows/graceful-shutdown.md) - On a stop signal the server cancels every live job, closes the listener within 5 seconds and exits with status 0.
* [Job lifecycle workflow](/workflows/job-lifecycle.md) - The multi-step process of a background computation: validate, acquire a slot, register, accept, compute with progress, terminate, release, expire.
* [Quality gate workflow](/workflows/quality-gate.md) - The test-first verification sequence: format, vet, race-enabled tests, frontend syntax, build, live probes, contrast check — ending in a named green state.
* [Sudoku solve / count workflow](/workflows/sudoku-solve.md) - Mask-guided, MRV-ordered backtracking used both to count solutions up to a cap and to find the first solution of a sudoku grid.
