# Getting started

## Install

```sh
go get github.com/JoddabodScripts/neri-go
```

Requires Go 1.21+ (generics are used internally for the LRU cache).

## Create a bot token

Create a bot on Nerimity and grab its token. Treat it like a password — anyone
with it can act as your bot.

## A minimal bot

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
			return // ignore our own messages
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

Run it with `go run .`. `Login` blocks and drives the connection until the
process exits, `client.Close()` is called, or the context passed to
`LoginWithContext` is cancelled.

## Options

`nerimity.New` takes an `Options` struct. Every field is optional; the zero
value connects to production Nerimity with a 1000-message cache:

```go
client := nerimity.New(nerimity.Options{
	// Point at a self-hosted instance instead of nerimity.com.
	WSURLOverride:  "https://my-instance.example.com",
	APIURLOverride: "https://my-instance.example.com",

	// LRU cache size for Client.messages. -1 for unbounded, 0 (default) for 1000.
	MessageCacheLimit: 5000,

	// Bring your own *http.Client (custom transport, proxy, timeouts, ...).
	HTTPClient: myHTTPClient,
})
```

## Connecting and reconnecting

`Login` authenticates over the Nerimity WebSocket gateway. If the connection
drops, the client automatically reconnects with exponential backoff (1s
doubling up to a 5s cap, with jitter — the same timing the JS SDK's underlying
`socket.io-client` uses) and re-runs the handshake, which repopulates every
cache and fires `OnReady` again. Handlers don't need to do anything special to
survive a reconnect; just don't assume `OnReady` fires exactly once.

## Where to go next

- [`events.md`](events.md) for the full event list
- [`messages.md`](messages.md) for sending, replying, editing, buttons,
  attachments, and HTML embeds
- [`permissions.md`](permissions.md) for the permission bit reference
- [`slash-commands.md`](slash-commands.md) for registering commands and
  routing sub-commands
- [`mentions.md`](mentions.md) for mention parsing
- [`webhooks.md`](webhooks.md) for sending messages without a bot account
