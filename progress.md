# Progress

## First loop
`sporecaster verify` — load spore.yaml, run Go structural checks, print a pass/fail report.
`sporecaster inspect` — serve a local web UI visualising packages, imports, exports, and channels.

## Completed

### Spec
- All packages declared with names, kinds, descriptions, imports, exports
- `error_handling` section (terminates_at, sentinels_are_exported)
- `channels` section with CLI commands (verify, init, inspect)
- `visualizer` and `visualizer/server` sub-package added and specified

### Verify checks
- Zoom 1 — skeleton: package directories exist, declared files present, Go package declarations correct
- Zoom 2 — connections: import graph (declared vs actual), CLI commands handled in main, os.Exit only in main, sentinel errors exported
- Zoom 3 — exports: exported symbols match spec by name and kind (struct / interface / func); struct fields verified by name and type; func signatures verified (params + return types); interface methods verified by name and signature — both directions (missing and undeclared)

### Spec
- All packages declared with names, kinds, descriptions, imports, exports
- `error_handling` section (terminates_at, sentinels_are_exported)
- `channels` section with CLI commands (verify, init, inspect)
- `visualizer` and `visualizer/server` sub-package added and specified
- All struct exports annotated with fields (name + type)
- All func exports annotated with signatures
- All interface exports annotated with methods and their signatures
- `Field`, `Method` types added to spore package to support spec completeness

### Visualiser (`sporecaster inspect`)
- Package cards with kind badges, import badges, click-to-expand exports
- Struct exports show fields (name + type) with left-border indent
- Func exports show signature in monospace below the export name
- Interface exports show methods with signatures (purple left-border indent)
- Sub-packages nested inside parent cards
- Channels section with commands and usage strings

### Agents
- Base: architect.md, builder.md, debugger.md (stubs)
- Golang: architect.md, builder.md (stubs)

---

## Next loops

### Loop 2 — goroutines and channels
- Add `routines:` section to spore.yaml per package (name, trigger, owns channel)
- Add `channels:` per package (name, type, direction) — distinct from top-level CLI/HTTP channels
- Zoom 2 check: verify goroutine spawns and channel declarations exist in code
- Visualiser: show routines and channels on package cards
- Implement concurrent check execution in `verify/golang`

### Loop 3 — tests
- Add `tests:` section to spore.yaml per package (kind: unit/integration, covers: [...])
- Zoom check: verify `_test.go` files exist for packages that declare tests
- Verify test function names match declared coverage targets
- Visualiser: show test coverage status on package cards

### Loop 4 — user flows
- Add `flows:` section to spore.yaml (name, steps: actor + action)
- Steps reference packages by name, making flows machine-verifiable
- Zoom check: trace that the actors named in each step are wired correctly (complements import graph)
- Scaffolding: generate integration test skeletons from declared flows
- Visualiser: flow diagram tab showing step sequences

### Loop 5 — agent system
- **Coder agent**: takes a package spec (name, description, exports, imports) and implements it
- **Checker agent**: runs `sporecaster verify` after coder acts, feeds failures back in a loop
- **Debugger agent**: reads a FAIL report, traces the call stack through the spec, identifies responsible package
- **Architect agent**: rewrite to be zoom-aware — operates at the right abstraction level per session
- **Test writer agent**: generates `_test.go` stubs from declared exports and flows

### Loop 6 — HTTP channel and web scale
- Extend `channels:` to support `type: http` with `routes:` (method, path, handler package)
- `type: websocket` channel for real-time
- Frontend components as packages with props/emits as exports
- Zoom check: verify route handlers exist and are wired to the declared handler package
- Visualiser: route table view for HTTP channels

---

## Known gaps (not yet in any loop)

- Exported `var`/`const` not covered by Zoom 3 (sentinel errors will surface here)
- External dependencies (`gopkg.in/yaml.v3` etc.) not declared in spec
- `verify/rails` — interface satisfied, no checks implemented
- Interface satisfaction not structurally verified (e.g. `verify/golang.Checker` implements `verify.Checker`)
- `sporecaster init` command not implemented

## Deferred
- `visualizer` HTTP server goroutine — currently blocks in `http.ListenAndServe`; needs graceful shutdown and a server goroutine declaration in the spec once Loop 2 lands
