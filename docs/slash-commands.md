# Slash commands

## Registering commands

```go
err := client.UpdateCommands(ctx, "your bot token", []nerimity.CommandDefinition{
	{Name: "help", Description: "Shows the help list", Args: "<page number>"},
	{Name: "ping", Description: "Checks if the bot is alive"},
})
```

`UpdateCommands` replaces the bot's entire command list; pass every command
you want to keep, not just the ones changing. It authenticates with the token
you pass directly (not the session from `Login`), so it can be called without
an active connection; typically you run this once, out of band, whenever your
command set changes, not on every bot startup.

Each `Description` is capped at 60 characters by the Nerimity server.
`UpdateCommands` checks this client-side and returns an error naming the
offending command before making any request, so you find out immediately
instead of from an opaque API error:

```go
err := client.UpdateCommands(ctx, token, []nerimity.CommandDefinition{
	{Name: "help", Description: "This description is deliberately way too long to pass validation"},
})
// err: nerimity: command "help" description is 68 characters, exceeds the 60-character limit
```

## Handling a command

A message matching `/name:<yourBotID> args...` gets a parsed `Command` on
`Message.Command`:

```go
client.OnMessageCreate(func(m *nerimity.Message) {
	if m.Command == nil {
		return
	}
	switch m.Command.Name {
	case "help":
		page := "1"
		if len(m.Command.Args) > 0 {
			page = m.Command.Args[0]
		}
		m.Reply(context.Background(), "Help page "+page)
	case "ping":
		m.Reply(context.Background(), "Pong!")
	}
})
```

`m.Command` is only set when the message targets this bot's user ID; commands
sent to other bots don't parse as yours, even if the same channel sees both.

## Sub-command routing pattern

The SDK deliberately does not include a command framework (matching the JS
SDK's philosophy: this is a client library, not a bot framework). For commands
with sub-verbs, the recommended pattern is to treat the first argument as the
verb and dispatch on it yourself:

```go
client.OnMessageCreate(func(m *nerimity.Message) {
	if m.Command == nil || m.Command.Name != "role" {
		return
	}
	args := m.Command.Args
	if len(args) == 0 {
		m.Reply(context.Background(), "usage: /role:bot <add|remove|list> ...")
		return
	}

	verb, rest := args[0], args[1:]
	switch verb {
	case "add":
		handleRoleAdd(m, rest)
	case "remove":
		handleRoleRemove(m, rest)
	case "list":
		handleRoleList(m)
	default:
		m.Reply(context.Background(), fmt.Sprintf("unknown subcommand %q", verb))
	}
})
```

This scales to a small router function per top-level command if you have many
commands; a `map[string]func(*nerimity.Message, []string)` keyed by command
name is a natural next step, and you're free to build one; the SDK just
doesn't ship one for you.
