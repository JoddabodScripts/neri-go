// Package nerimity is a client SDK for the Nerimity chat platform
// (https://nerimity.com), functionally equivalent to the official JavaScript
// SDK, @nerimity/nerimity.js.
//
// # Quickstart
//
//	client := nerimity.New(nerimity.Options{})
//
//	client.OnReady(func() {
//		log.Printf("connected as %s", client.User().Username)
//	})
//
//	client.OnMessageCreate(func(m *nerimity.Message) {
//		if m.User.ID == client.User().ID {
//			return
//		}
//		if m.Content == "!ping" {
//			m.Reply(context.Background(), "Pong!")
//		}
//	})
//
//	if err := client.Login("bot token here"); err != nil {
//		log.Fatal(err)
//	}
//
// # Events
//
// Handlers are registered with typed On* methods on Client (OnReady,
// OnMessageCreate, OnMessageButtonClick, and so on) rather than a single
// stringly-typed dispatch method. This keeps handler signatures discoverable
// via godoc and type-checked at compile time. Multiple handlers may be
// registered for the same event; they run in registration order on a single
// internal dispatch goroutine, so handler code does not need its own
// synchronization to safely read the client's caches.
//
// Login blocks the calling goroutine, running the connection (with automatic
// reconnect) until Close is called or its context is cancelled. Run it in your
// program's main goroutine, or in its own goroutine if you need main for
// something else.
//
// See docs/events.md for the full event list and docs/getting-started.md for
// a longer walkthrough.
package nerimity
