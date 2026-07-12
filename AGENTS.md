# AGENTS.md

Instructions for coding agents working in this repository.

## Project layout

This is a flat Go package at the module root: `package nerimity`, module path
`github.com/JoddabodScripts/neri-go`. There is no internal/ or pkg/ split — the
surface area doesn't justify one. Files are organized by concept, not by
layer:

- `client.go` — `Client`, `Options`, connection lifecycle (`Login`, `Close`),
  reconnect/backoff, the dispatch goroutine, cache stores
- `events.go` — `On*` handler registration and the WebSocket event → cache
  update → handler dispatch logic
- `socketio.go` — the minimal Engine.IO v4 / Socket.IO v5 client over
  `gorilla/websocket`. This is the only file that speaks the wire protocol at
  the frame level
- `rawdata.go` — unexported structs (`raw*`) that mirror the JSON the server
  actually sends, used only for unmarshalling
- `rest.go` — REST calls (send/edit/delete message, ban/kick, button callback)
- `message.go`, `channel.go`, `user.go`, `member.go`, `server.go`, `role.go`,
  `reaction.go`, `button.go` — the exported domain types
- `mentions.go`, `permissions.go`, `html.go` — pure parsing/helper logic with
  no network dependency
- `attachment.go`, `webhook.go`, `commands.go` — attachment upload, webhook
  sending, slash-command registration
- `cache.go` — the generic LRU cache backing every collection (`Client.Servers`,
  `Server.Members`, the message cache, etc.)
- `examples/` — runnable example bots, also serves as compiled documentation
- `docs/` — prose documentation

## Running checks

```sh
go build ./...
go vet ./...
gofmt -l .          # should print nothing; gofmt -w . to fix
go test ./...
```

All four must pass before a change is done. `gofmt -l .` printing any path is
a failure, not a suggestion.

## Wire protocol verification — read this before touching client.go, events.go,
## socketio.go, rawdata.go, or rest.go

Nerimity has no public API spec. The only source of truth for event names,
REST endpoints, and payload shapes is the reference JavaScript SDK at
`~/nrepos/nerimity.js` (`src/classes/Client.ts`, `src/EventNames.ts`,
`src/RawData.ts`, `src/services/`). If that path doesn't exist in your
environment, ask for it before guessing — do not invent field names or
endpoints from training data. When adding a new event or endpoint:

1. Find the equivalent in nerimity.js first — the exact socket event string
   (`src/EventNames.ts`), the payload shape (`src/RawData.ts`), and how the JS
   `Client` mutates its caches in response (`src/classes/Client.ts`'s
   `EventHandlers` class).
2. Mirror the cache mutation logic, not just the event name. The JS SDK's
   `onXxx` handlers in `EventHandlers` show exactly what gets added, removed,
   or patched in which collection — get this wrong and consumers see stale or
   missing data even though the event fired.
3. Add the raw struct to `rawdata.go`, the public type/field to the relevant
   domain file, the socket event constant and handler to `events.go`.

## Conventions

- **Errors**: return `(T, error)`, never panic on request failures. Use
  `fmt.Errorf("nerimity: doing the thing: %w", err)` — always prefix with
  `nerimity: ` so errors are identifiable when this package is one dependency
  among many in a larger bot.
- **Context**: every method that makes a network call takes `context.Context`
  as its first argument.
- **Naming**: exported types/methods use Go naming (`ServerID`, not `serverId`);
  unexported `raw*` structs keep the wire's camelCase via `json:` tags, not in
  the Go field names.
- **Caches**: never expose the internal `cache[T]` type directly. Public
  accessors return `[]T` snapshots (see `Server.Members()`,
  `Client.Servers()`) so callers can't corrupt cache internals or deadlock by
  holding a lock across a callback.
- **New events**: add the `On<EventName>` registration method in `events.go`
  next to the others, add the socket event name constant, add a `case` in
  `Client.handleEvent`, and update `docs/events.md`.
- **Tests**: anything with parsing/pure logic (mention parsing, command
  parsing, HTML escaping, permission bit helpers, the LRU cache, Socket.IO
  frame parsing) must have unit tests — they don't need a live server. Files
  needing a live Nerimity connection are integration-only and are not part of
  `go test ./...`; there are none currently, and none should be added without
  a way to skip them in CI.
