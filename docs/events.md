# Events

Every event is registered with a typed `On<Name>` method on `Client`. You can
register more than one handler per event; they run in registration order on a
single internal goroutine, so they can safely read the client's caches without
extra locking. Don't block for a long time inside a handler — it delays every
other event until it returns.

## Ready

```go
client.OnReady(func() {
	log.Printf("connected as %s", client.User().Username)
})
```

Fires once the client has authenticated and populated its caches (servers,
channels, members, roles). Fires again after every automatic reconnect.

## MessageCreate

```go
client.OnMessageCreate(func(m *nerimity.Message) {
	fmt.Println(m.User.Username, "said", m.Content)
})
```

Fires when any message is created in a channel the bot can see. See
[`messages.md`](messages.md) for the `Message` type.

## MessageUpdate

```go
client.OnMessageUpdate(func(m *nerimity.Message) {
	fmt.Println("edited:", m.Content)
})
```

Fires when a message is edited. `m` is the same cached `*Message`, mutated in
place, so any earlier references you kept to it will also reflect the edit.

## MessageDelete

```go
client.OnMessageDelete(func(e nerimity.MessageDeleteEvent) {
	fmt.Println("deleted message", e.MessageID, "in channel", e.ChannelID)
})
```

Fires when a message is deleted. Only IDs are available — by the time this
fires, the message has already been evicted from the cache.

## MessageReactionAdded / MessageReactionRemoved

```go
client.OnMessageReactionAdded(func(r *nerimity.Reaction, user *nerimity.User) {
	// user is nil if the reacting user isn't cached.
	fmt.Printf("%v reacted %s (now %d)\n", user, r.Name, r.Count)
})

client.OnMessageReactionRemoved(func(r *nerimity.Reaction, user *nerimity.User) {
	fmt.Printf("%v removed %s (now %d)\n", user, r.Name, r.Count)
})
```

`r.Message` may be nil if the message wasn't cached; check `r.Partial`.

## MessageButtonClick

```go
client.OnMessageButtonClick(func(b *nerimity.MessageButton) {
	if b.ID != "confirm" {
		return
	}
	b.Respond(context.Background(), nerimity.ButtonResponse{
		Content: "Confirmed!",
	})
})
```

Fires when a user clicks a button you attached to a message. See
[`messages.md`](messages.md#buttons) for sending buttons and responding.

## ServerMemberJoined / ServerMemberLeft / ServerMemberUpdated

```go
client.OnServerMemberJoined(func(m *nerimity.ServerMember) {
	fmt.Println(m.User.Username, "joined", m.Server().Name)
})

client.OnServerMemberLeft(func(m *nerimity.ServerMember) {
	fmt.Println(m.User.Username, "left")
})

client.OnServerMemberUpdated(func(m *nerimity.ServerMember) {
	fmt.Println(m.User.Username, "now has roles", m.RoleIDs)
})
```

`ServerMemberUpdated` currently reports role changes (`m.RoleIDs`), matching
what the gateway sends.

## ServerJoined / ServerLeft

```go
client.OnServerJoined(func(s *nerimity.Server) {
	fmt.Println("joined server:", s.Name)
})

client.OnServerLeft(func(s *nerimity.Server) {
	fmt.Println("left server:", s.Name)
})
```

Fires when the bot itself joins or leaves (or is kicked/banned from) a server.
`ServerLeft` fires before the server and its channels/members are purged from
cache, so `s` and its collections are still fully populated inside the
handler.

## ServerChannelCreated / ServerChannelUpdated / ServerChannelDeleted

```go
client.OnServerChannelCreated(func(c *nerimity.ServerChannel) {
	fmt.Println("new channel:", c.Name)
})

client.OnServerChannelUpdated(func(c *nerimity.ServerChannel) {
	fmt.Println("channel renamed to:", c.Name)
})

client.OnServerChannelDeleted(func(e nerimity.ServerChannelDeleteEvent) {
	fmt.Println("deleted channel", e.ChannelID, "in server", e.ServerID)
})
```

`ServerChannel` is an alias for `Channel` — see [`messages.md`](messages.md)
for its fields.

## ServerRoleCreated / ServerRoleUpdated / ServerRoleDeleted

```go
client.OnServerRoleCreated(func(r *nerimity.ServerRole) {
	fmt.Println("new role:", r.Name)
})

client.OnServerRoleUpdated(func(r *nerimity.ServerRole) {
	fmt.Println("role updated:", r.Name)
})

client.OnServerRoleDeleted(func(r *nerimity.ServerRole) {
	fmt.Println("role deleted:", r.Name)
})
```

## ServerRoleOrderUpdated

```go
client.OnServerRoleOrderUpdated(func(server *nerimity.Server, roles []*nerimity.ServerRole) {
	for _, r := range roles {
		fmt.Println(r.Order, r.Name)
	}
})
```

Fires when a server's role hierarchy is reordered. `roles` is every cached
role for that server, already updated to the new `Order` values.

## Error

```go
client.OnError(func(err error) {
	log.Printf("nerimity: dropped a gateway event: %v", err)
})
```

Fires whenever the client fails to decode a gateway event's payload — normally
because the server sent a shape this SDK version doesn't recognize for that
event. The event in question is dropped: its own handlers did not run, and any
cache updates it would have made did not happen.

Registering `OnError` is optional but strongly recommended in production. The
one case worth calling out specifically: if part of the `Ready`/authentication
payload fails to decode, the client still sets `Client.User()` and fires
`OnReady` as long as the `user` field itself parsed correctly — only the
malformed sub-field (say, `serverRoles`) is skipped and reported here. Without
an `OnError` handler you'd have no visibility into that partial failure at
all.
