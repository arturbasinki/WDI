# domains

* [Background computation (domain)](/domains/background-computation.md) - The supporting domain that turns long-running computations into identified, pollable, cancellable, resource-bounded background work.
* [Browser presentation (domain)](/domains/browser-presentation.md) - The browser-side presentation domain: one hash-routed single-page client rendering all eight task views without a build step.
* [Cryptarithm puzzles (z07)](/domains/cryptarithm.md) - Solve letter-for-digit equation puzzles over + - * / by pruned permutation search, honouring digit bijectivity, no leading zeros, and exact division, and report every solution.
* [Function plotting (z06)](/domains/function-plotting.md) - Parse a learner-written function expression and sample it over an interval, detecting vertical asymptotes numerically and returning the curve as segments plus a clamp threshold.
* [Linear systems (z01)](/domains/linear-systems.md) - Solve N×N systems of linear equations by pivoted Gauss-Jordan reduction, classifying unique, contradictory and underdetermined cases, and record the reduction steps.
* [Narcissistic numbers (z02)](/domains/narcissistic-numbers.md) - Exhaustive, parallel, incrementally optimised search for numbers equal to the sum of their digits raised to the digit-count power, to a verified reach of 10^12.
* [Gronsfeld cipher (z05)](/domains/polyalphabetic-cipher.md) - Gronsfeld polyalphabetic cipher with a numeric key over a 27-symbol alphabet, its inverse, and a dictionary-attack breaker over enumerated keys.
* [Prime sieving (z03)](/domains/prime-sieving.md) - Sieve of Eratosthenes over three data structures (boolean slice, map, bitmap) with a cancellable, sequential timing comparison across increasing ranges.
* [Sudoku generation and solving (z08)](/domains/sudoku.md) - Generate sudoku puzzles with a provably unique solution and many empty cells, by random fill plus uniqueness-preserving dig, and solve/count any grid.
* [Task catalog (domain)](/domains/task-catalog.md) - The single authoritative registry of the eight tasks (identity, number, title, blurb) that the whole presentation layer reads.
* [Text word frequency (z04)](/domains/text-frequency.md) - Case-insensitive word-frequency ranking over a text with a five-letter floor, canonically the top 20 words of an embedded literary text.
