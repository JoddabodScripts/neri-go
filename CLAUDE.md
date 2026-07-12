# CLAUDE.md

This file is for Claude Code specifically. The full set of conventions, project
layout, and wire-protocol verification rules lives in
[`AGENTS.md`](AGENTS.md) — read that first, it's the source of truth.

A few things worth calling out explicitly for Claude:

- **Don't guess at the Nerimity wire protocol.** If you're touching
  `client.go`, `events.go`, `socketio.go`, `rawdata.go`, or `rest.go`, go read
  the matching code in `~/nrepos/nerimity.js` before writing anything. This
  codebase has no public API docs; that repo is the only ground truth. If it's
  not present in your environment, say so and ask rather than inventing an
  event name or payload shape.
- **Run the full check sequence before calling anything done**: `go build
  ./...`, `go vet ./...`, `gofmt -l .` (must print nothing), `go test ./...`.
  Don't report success on partial evidence.
- **This is a library, not an app.** There's no server to start, no UI to
  click through. "Testing the feature" means writing or running a `go test`
  for the parsing/logic pieces (mentions, commands, HTML escaping, permission
  bits, the LRU cache, Socket.IO frame parsing), and checking that
  `examples/*/main.go` still compiles against the current API — those examples
  are compiled documentation, not throwaway scratch files.
- **Keep the flat package.** Don't split `nerimity` into subpackages
  (`internal/`, `pkg/`) unless asked — the surface area doesn't warrant it, and
  a newcomer's `import "github.com/JoddabodScripts/neri-go"` should stay a
  single import for the whole SDK.
