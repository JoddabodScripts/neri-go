# neri-go

A Go SDK for [Nerimity](https://nerimity.com), functionally equivalent to the
official JavaScript SDK, [`@nerimity/nerimity.js`](https://github.com/Nerimity/nerimity.js).

## Install

```sh
go get github.com/JoddabodScripts/neri-go
```

## Quickstart

```go
package main

import (
	"context"
	"log"

	"github.com/JoddabodScripts/neri-go"
)

func main() {
	client := nerimity.New(nerimity.Options{})

	client.OnReady(func() {
		log.Printf("connected as %s", client.User().Username)
	})

	client.OnMessageCreate(func(m *nerimity.Message) {
		if m.User == nil || m.User.ID == client.User().ID {
			return
		}
		if m.Content == "!ping" {
			m.Reply(context.Background(), "Pong!")
		}
	})

	if err := client.Login("your bot token"); err != nil {
		log.Fatal(err)
	}
}
```

That's the whole bot. `Login` blocks and drives the connection (with automatic
reconnect) until `client.Close()` is called or its context is cancelled.

## Event model

Handlers are registered with typed `On*` methods on `Client` — `OnReady`,
`OnMessageCreate`, `OnMessageButtonClick`, `OnServerMemberJoined`, and so on —
rather than a single stringly-typed emitter. This keeps every handler's
signature discoverable via `go doc` and checked by the compiler. You can
register more than one handler for the same event; they run, in registration
order, on a single internal dispatch goroutine, so handlers can safely read the
client's caches without their own locking.

See [`docs/events.md`](docs/events.md) for the full event list with payload
types and example handlers.

## Documentation

- [`docs/getting-started.md`](docs/getting-started.md) — a longer walkthrough
- [`docs/events.md`](docs/events.md) — every event, its payload, and an example
- [`docs/messages.md`](docs/messages.md) — sending, replying, editing, buttons,
  attachments, HTML embeds
- [`docs/webhooks.md`](docs/webhooks.md) — sending messages without a bot
- [`docs/permissions.md`](docs/permissions.md) — the permission bit reference
- [`docs/slash-commands.md`](docs/slash-commands.md) — registration and the
  sub-command routing pattern
- [`docs/mentions.md`](docs/mentions.md) — mention parsing and stripping

Worked examples live in [`examples/`](examples/): an echo bot, a
button-interaction bot, and a webhook sender.

## Design notes

- No built-in command framework, matching the JS SDK's philosophy. Sub-command
  routing (first token as the verb, the rest as arguments) is documented as a
  pattern in [`docs/slash-commands.md`](docs/slash-commands.md), not forced on
  you as a framework.
- Sends on a single channel are queued and executed in order, matching the JS
  SDK's per-channel `AsyncFunctionQueue` behavior — concurrent calls to
  `channel.Send` never race or reorder messages.
- The message cache is an LRU with a default limit of 1000, matching the JS
  SDK's `Collection`.

## Contributing

See [`AGENTS.md`](AGENTS.md) (or [`CLAUDE.md`](CLAUDE.md) if you're using
Claude Code) for project layout, test/lint commands, and conventions.

## License

MIT — see [`LICENSE`](LICENSE).
