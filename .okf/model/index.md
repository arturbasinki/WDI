# model

* [Cipher alphabet](/model/cipher-alphabet.md) - The 27-symbol cipher alphabet (space + 26 lowercase letters) whose ordering and size define the modular shift arithmetic.
* [Cipher key](/model/cipher-key.md) - A digit-string key whose digits supply the per-position shifts, cycling over its length; leading zeros are significant, length 1..9.
* [Cryptarithm puzzle](/model/cryptarithm-puzzle.md) - Aggregate of a cryptarithm: the equation, its letters and no-zero letters, the search space size, the full solution count and the (possibly capped) solution list.
* [Elimination step](/model/elimination-step.md) - A single narrated elimination operation (description + full matrix snapshot) used to animate the reduction, bounded by a cell budget.
* [Embedded corpus](/model/embedded-corpus.md) - Reference text shipped inside the deployed artifact — the canonical literary corpus and the English dictionary — identified by name and immutable at run time.
* [Expression](/model/expression.md) - Aggregate of a parsed arithmetic expression over x, constants and functions: precedence, associativity, implied multiplication, arity-checked calls, bounded depth.
* [Job](/model/job.md) - Aggregate for one background computation: opaque id, monotonic status lifecycle, progress/result snapshots, cancellation handle and end time, all lock-protected.
* [Linear system](/model/linear-system.md) - Aggregate of an augmented N×(N+1) matrix: square, exactly-supplied, finite coefficients, with rank/pivots/RREF and a parametric solution derived from it.
* [Normalization report](/model/normalization-report.md) - The outcome of preparing text for the cipher: the encipherable text plus the sorted set of unique characters dropped, so nothing disappears silently.
* [Plot clamping](/model/plot-clamping.md) - The IQR-derived ±Y threshold applied to outlier samples, applied by keeping them at the threshold so branches reach the frame of the plot.
* [Plot segment](/model/plot-segment.md) - A maximal continuous stretch of a plotted curve, broken at non-finite samples and asymptotic jumps so the renderer never draws across an asymptote.
* [Plot](/model/plot.md) - Aggregate of a plotted function: formula, interval, requested vs applied sample count, continuous segments, asymptote positions and the applied Y clamp.
* [Search progress](/model/search-progress.md) - The pollable state of a long computation: processed/limit, results so far as a sorted copy, and a done flag reserved for genuinely finished runs.
* [Sieve benchmark row](/model/sieve-benchmark-row.md) - One measurement of one sieve data structure over one range: structure, n, prime count and duration — the row of the comparison table the exercise asks for.
* [Solution status](/model/solution-status.md) - The three-way classification of a linear system (unique, contradictory, underdetermined) and the rendering rules that follow from each state.
* [Sudoku puzzle](/model/sudoku-puzzle.md) - Aggregate of a 81-cell sudoku position: generated puzzle with its proven-unique solution and blank count, plus a solve/count result for arbitrary grids.
* [Task descriptor](/model/task-descriptor.md) - Entity for one catalogued task: stable id (the exercise number), display number, title and blurb, served as the single source of task identity.
* [Word tally](/model/word-tally.md) - A word and its occurrence count in an analysed text, the unit of the frequency ranking with case folding and a length floor.
