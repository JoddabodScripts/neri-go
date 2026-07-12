# Mentions

Nerimity encodes mentions as raw tokens in message content:

| Token          | Meaning                          |
|----------------|-----------------------------------|
| `[@:<userId>]` | Mention of a specific user        |
| `[@:e]`        | `@everyone`                       |
| `[@:s]`        | `@someone` (random online member) |

## Parsed mentions

`Message.Mentions` is a `[]Mention`, parsed in order of appearance:

```go
type Mention struct {
	Type   MentionType // MentionTypeUser, MentionTypeEveryone, or MentionTypeSomeone
	UserID string      // set only when Type == MentionTypeUser
	Raw    string       // the exact token, e.g. "[@:8371002...]"
}
```

Matching the JS SDK, mentions are **not** resolved to display names here; you
get the raw token and, for user mentions, the ID. If you need the mentioning
user's `User` object, look it up yourself:

```go
for _, mention := range m.Mentions {
	if mention.Type != nerimity.MentionTypeUser {
		continue
	}
	if user := client.GetUser(mention.UserID); user != nil {
		fmt.Println("mentioned:", user.Username)
	}
}
```

## Mentioning someone in a message you send

Build the token yourself, or use the `Mention()` helper on `User` /
`ServerMember` / `Channel`:

```go
channel.Send(ctx, "Hey "+user.Mention()+", check this out")
// equivalent to: channel.Send(ctx, "Hey [@:"+user.ID+"], check this out")
```

## Stripping mentions

If your bot wants to process a message's plain text without mention markup:

```go
plain := nerimity.StripMentions(m.Content)
```

This removes every `[@:...]` token and collapses the whitespace left behind.
