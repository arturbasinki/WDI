# invariants

* [Cipher alphabet closure](/invariants/cipher-alphabet-closure.md) - Cipher input is confined to the 27-symbol alphabet; dropped characters are always reported and empty-after-normalization is an error.
* [Distinct digit assignment](/invariants/distinct-digit-assignment.md) - Cryptarithm solutions are digit bijections: unique digits per letter, at most ten letters, no leading zero on multi-digit numbers, exact division.
* [Expression depth bound](/invariants/expression-depth-bound.md) - Expression parsing is bounded to 100 nesting levels and 2000 characters so pathological input is a rejection, never a process crash.
* [Job status monotonicity](/invariants/job-status-monotonic.md) - A job's status is monotonic: running then exactly one terminal state, with the end timestamp written atomically with it.
* [Linear system size and finiteness](/invariants/linear-system-size.md) - A linear system has 1..100 equations, exactly N*(N+1) finite numeric values, with NaN and infinities rejected before solving.
* [Sieve structure agreement](/invariants/sieve-structure-agreement.md) - All three sieve structures yield the same ascending prime set for a given n, matching published reference counts.
* [Sudoku grid validation](/invariants/sudoku-grid-validation.md) - A sudoku grid is exactly 81 values each in 0..9, and solution counting requires a cap of at least 1.
* [Sudoku unique solution](/invariants/sudoku-unique-solution.md) - Every generated sudoku has exactly one solution, verified by counting to two after each cell removal.
* [Task catalog consistency](/invariants/task-catalog-consistency.md) - The task catalog has exactly eight uniquely-identified entries in order, all fields populated, no duplicated copy on the client, and no blurb that over-claims.
* [Word length floor](/invariants/word-length-floor.md) - Only letter runs of at least five characters (case-folded) enter the word-frequency tally.
