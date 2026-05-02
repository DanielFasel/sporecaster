# Progress

## First loop
`sporecaster verify` — load spore.yaml, run Go structural checks, print a pass/fail report.

## Completed
- All packages: Zoom 1 — skeleton (names, kinds, descriptions)
- All packages: Zoom 2 — connections (imports, error_handling, CLI channel)

## In progress
- Zoom 3: package-by-package exports, starting with `spore` (shared types)

## Open questions
(none)

## Deferred
- `visualizer` package — web UI for exploring a spore; will need an HTTP channel
  and a server goroutine when zoomed into
- `verify/rails` implementation — interface is satisfied, no checks yet
- Concurrent check execution in `verify/golang` — no goroutines needed for first loop
