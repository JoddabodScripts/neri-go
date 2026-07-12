package nerimity

// Reaction is an emoji reaction on a message, delivered with the
// MessageReactionAdded and MessageReactionRemoved events.
type Reaction struct {
	client *Client

	// Name is the emoji name (or the unicode emoji itself for standard
	// emojis).
	Name string
	// EmojiID is the custom emoji's ID, empty for standard emojis.
	EmojiID string
	// Gif reports whether the custom emoji is animated.
	Gif bool
	// Count is the number of reactions with this emoji after the event.
	Count int
	// MessageID is the ID of the reacted-to message.
	MessageID string
	// ChannelID is the ID of the channel the message is in.
	ChannelID string
	// Message is the reacted-to message if it is cached, else nil (see
	// Partial).
	Message *Message
	// Channel is the channel the message is in if cached, else nil.
	Channel *Channel
	// Partial reports whether Message could not be resolved from cache.
	Partial bool
}

func newReactionFromAdded(client *Client, p reactionAddedPayload) *Reaction {
	return newReaction(client, p.MessageID, p.ChannelID, p.Name, p.EmojiID, p.Count, p.Gif)
}

func newReactionFromRemoved(client *Client, p reactionRemovedPayload) *Reaction {
	return newReaction(client, p.MessageID, p.ChannelID, p.Name, p.EmojiID, p.Count, p.Gif)
}

func newReaction(client *Client, messageID, channelID, name, emojiID string, count int, gif bool) *Reaction {
	r := &Reaction{
		client:    client,
		Name:      name,
		EmojiID:   emojiID,
		Gif:       gif,
		Count:     count,
		MessageID: messageID,
		ChannelID: channelID,
		Partial:   true,
	}
	if msg, ok := client.messages.get(messageID); ok {
		r.Message = msg
		r.Partial = false
	}
	r.Channel, _ = client.channels.get(channelID)
	return r
}
