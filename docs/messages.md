# Messages

## The Message type

`Message` carries the resolved author and channel, parsed mentions, and,
for slash-command messages targeting this bot, a parsed `Command`:

```go
type Message struct {
	ID        string
	Content   string
	Type      MessageType
	ChannelID string
	Channel   *Channel
	User      *User
	CreatedAt int64
	EditedAt  int64
	Mentions  []Mention
	Command   *Command
	Replies   map[string]*Message
	WebhookID string
	// ...
}
```

`m.Member()` resolves the author's `*ServerMember` for the channel's server
(nil for DMs).

## Sending

```go
msg, err := channel.Send(ctx, "hello!")
```

`Channel.Send` takes an optional `MessageOptions`:

```go
type MessageOptions struct {
	HTMLEmbed         string
	NerimityCdnFileID string
	Buttons           []ButtonOption
	Silent            bool
	ReplyToMessageIDs []string
	MentionReplies    bool
}
```

Sends on the same channel are serialized (a per-channel queue), so calling
`Send` concurrently from multiple goroutines never reorders or interleaves
requests — the same guarantee the JS SDK's `AsyncFunctionQueue` gives you.

## Replying

```go
client.OnMessageCreate(func(m *nerimity.Message) {
	if m.Content == "!ping" {
		m.Reply(context.Background(), "Pong!")
	}
})
```

`Message.Reply` is `Channel.Send` with `ReplyToMessageIDs` set to this
message's ID.

## Editing

```go
edited, err := msg.Edit(ctx, "updated content")
```

Only the bot's own messages can be edited. `EditOptions` lets you also replace
the HTML embed or buttons.

## Deleting

```go
err := msg.Delete(ctx)
// or, from a channel, by ID:
err := channel.DeleteMessage(ctx, messageID)
```

## Buttons

Attach buttons when sending:

```go
channel.Send(ctx, "Pick one:", nerimity.MessageOptions{
	Buttons: []nerimity.ButtonOption{
		{ID: "yes", Label: "Yes"},
		{ID: "no", Label: "No"},
	},
})
```

Handle clicks with `OnMessageButtonClick`, and respond with content, a title,
a button-label override, and/or interactive components:

```go
client.OnMessageButtonClick(func(b *nerimity.MessageButton) {
	switch b.ID {
	case "yes":
		b.Respond(context.Background(), nerimity.ButtonResponse{Content: "Confirmed."})
	case "no":
		b.Respond(context.Background(), nerimity.ButtonResponse{
			Title:   "Tell us more",
			Content: "Why not?",
			Components: []nerimity.ButtonComponent{
				nerimity.InputComponent("reason", "Reason", "Type here..."),
			},
		})
	}
})
```

Component builders: `nerimity.TextComponent`, `nerimity.DropdownComponent`,
`nerimity.InputComponent`. Submitted values come back on a later
`MessageButtonClick` event (`Type == "modal_click"`) in `b.Data`, keyed by
component ID.

If the clicked message wasn't in the local cache, `b.Message` is nil and
`b.Partial` is true; call `b.Fetch(ctx)` to fetch it from the API.

## Attachments

Upload a file, then reference the returned CDN file ID when sending:

```go
attachment, err := nerimity.NewAttachmentFromFile("screenshot.png")
if err != nil {
	log.Fatal(err)
}
fileID, err := attachment.Build(ctx, client, message.Channel)
if err != nil {
	log.Fatal(err)
}
message.Reply(ctx, "here you go", nerimity.MessageOptions{
	NerimityCdnFileID: fileID,
})
```

Or build from any `io.Reader` with `nerimity.NewAttachment(r, "name.ext")`.

## HTML embeds

`MessageOptions.HTMLEmbed` renders an HTML block below the message. Nerimity's
server-side validator is stricter than a browser:

- Only an allowlisted set of tags and attributes is accepted.
- **Void elements must be self-closed**: write `<br/>` and `<img src="..."/>`,
  not `<br>` or `<img src="...">`.
- `position: fixed` is rejected outright.
- The validator counts opening and closing tags with a regular expression to
  check the markup is balanced. **Any `<`, `>`, `&`, `"`, or `'` character
  inside text content — including anything you didn't write yourself, like a
  username or message content you're echoing back — throws that count off and
  gets the whole embed rejected**, not just sanitized.

Because of that last point, escape any string you didn't write literally in Go
source before it goes into an `HTMLEmbed`:

```go
safeName := nerimity.EscapeHTML(user.Username)
embed := fmt.Sprintf(`<div>Welcome, %s!</div>`, safeName)
channel.Send(ctx, "", nerimity.MessageOptions{HTMLEmbed: embed})
```

`EscapeHTML` escapes `& < > " '`, in that order (ampersand first, so the
other substitutions' output doesn't get double-escaped).

## Mentions

See [`mentions.md`](mentions.md).
