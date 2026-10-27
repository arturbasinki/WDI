---
type: Decision
title: Content embedded in the artifact
description: "Compile the reference texts, the dictionary and the frontend into the deployed artifact, so canonical answers need no data files, paths or configuration."
tags: [decision, adr, deployment, content]
timestamp: "2026-10-27T16:32:44Z"
---

# Context

Two exercises are defined over specific bodies of text: the word-frequency ranking over
"Pana Tadeusz", the cipher attack over an English dictionary. A conventional layout puts them
in a data directory pointed at by configuration — which makes the canonical answers depend on
files being present, copied, and found.

# Decision

Compile the reference content **into the artifact**. Each task module owns the content it
analyses; there is no data directory, no path, no configuration. The frontend is embedded the
same way.

# Consequences

**Positive.** The default answers work on any machine with no preparation, which is what makes
"run the canonical example" a one-request demonstration. The container runtime stage carries
only the executable and builds with no network and no package installation. Content cannot be
replaced, traversed, or downloaded at run time.

**Negative.** The binary carries ~1 MB of text; changing content means rebuilding. A ~705 KB
duplication remains in the repository's legacy task folders, accepted deliberately as teaching
material. Duplicated content also lengthens compile time slightly.

**Confirmed by.** The running container needs one file; the image is ~15 MB; both canonical
answers are reproducible from a clean checkout.

# Related

[Embedded corpus](/model/embedded-corpus.md) ·
[Standard-library-only, single binary](/decisions/stdlib-only-single-binary.md) ·
[Static asset delivery](/contracts/static-assets.md) ·
[Corpus to ranking flow](/dataflows/corpus-to-ranking-flow.md)

# Citations

- `internal/z04/words.go`, `internal/z04/pan-tadeusz.txt` — the embedded literary text.
- `internal/z05/gronsfeld.go`, `internal/z05/slownik.txt` — the embedded dictionary.
- `cmd/wdi-server/main.go` — the embedded frontend filesystem.
- `Dockerfile` — the runtime stage carrying only the executable.
