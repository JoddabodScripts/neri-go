package nerimity

import (
	"context"
	"regexp"
	"strings"
)

// commandRegex matches a slash-command message of the form
// "/commandName:botUserId args...". Group 1 is "/name", group 2 is the bot
// user ID it targets, group 3 is the (optional) argument string.
var commandRegex = regexp.MustCompile(`(?m)^(/[^:\s]*):(\d+)( .*)?$`)

// Command is a parsed slash command from a message whose content matches
// "/name:botUserId args..." and targets this bot.
type Command struct {
	// Name is the command name without the leading slash.
	Name string
	// Args are the whitespace-separated arguments after the command token.
	Args []string
}

// Message is a chat message. It carries the resolved User and Channel where
// available, parsed mentions, and, for slash-command messages targeting this
// bot, a parsed Command.
type Message struct {
	client *Client

	// ID is the message's unique ID.
	ID string
	// Content is the message text, empty for contentless messages.
	Content string
	// Type is the message kind (content vs. a system event).
	Type MessageType
	// ChannelID is the ID of the channel the message was sent in.
	ChannelID string
	// Channel is the resolved channel, or nil if not cached.
	Channel *Channel
	// User is the message author, or nil if not cached.
	User *User
	// CreatedAt is the send time as a Unix millisecond timestamp.
	CreatedAt int64
	// EditedAt is the last-edit time as a Unix millisecond timestamp, 0 if
	// never edited.
	EditedAt int64
	// Mentions are the mention tokens parsed from Content, in order.
	Mentions []Mention
	// Command is the parsed slash command if Content matched
	// "/name:<thisBotID> ..."; otherwise nil.
	Command *Command
	// Replies are the messages this message was a reply to, keyed by ID.
	Replies map[string]*Message
	// WebhookID is set if the message was posted by a webhook.
	WebhookID string
}

func newMessage(client *Client, raw rawMessage) *Message {
	m := &Message{
		client:    client,
		ID:        raw.ID,
		Content:   raw.Content,
		Type:      raw.Type,
		ChannelID: raw.ChannelID,
		CreatedAt: raw.CreatedAt,
		EditedAt:  raw.EditedAt,
		WebhookID: raw.WebhookID,
		Replies:   map[string]*Message{},
	}

	m.Channel, _ = client.channels.get(raw.ChannelID)
	if m.Channel == nil {
		// Unknown channel (e.g. a DM we have not cached): create a minimal DM
		// channel so Reply/Send still work.
		m.Channel = client.channels.setChannel(rawChannel{ID: raw.ChannelID, Type: ChannelTypeDMText})
	}

	m.User, _ = client.users.get(raw.CreatedBy.ID)
	if m.User == nil {
		m.User = client.users.setUser(raw.CreatedBy)
	}

	for _, reply := range raw.ReplyMessages {
		if reply.ReplyToMessage != nil {
			rm := newMessage(client, *reply.ReplyToMessage)
			m.Replies[rm.ID] = rm
		}
	}

	if m.Content != "" {
		m.Mentions = parseMentions(m.Content)
		m.Command = parseCommand(m.Content, client.selfID())
	}
	return m
}

// parseCommand returns the parsed Command if content is a slash command
// targeting botUserID, or nil otherwise.
func parseCommand(content, botUserID string) *Command {
	if botUserID == "" {
		return nil
	}
	match := commandRegex.FindStringSubmatch(content)
	if match == nil || match[2] != botUserID {
		return nil
	}
	// Mirror the JS SDK exactly: args are content split on single spaces with
	// the command token dropped.
	parts := strings.Split(content, " ")
	return &Command{
		Name: strings.TrimPrefix(match[1], "/"),
		Args: parts[1:],
	}
}

// Member returns the message author's ServerMember for the channel's server, or
// nil for DMs or when not cached.
func (m *Message) Member() *ServerMember {
	if m.Channel == nil || m.User == nil {
		return nil
	}
	s := m.Channel.Server()
	if s == nil {
		return nil
	}
	return s.Member(m.User.ID)
}

// Reply sends a message that replies to this one. opts is optional; any
// ReplyToMessageIDs it sets are overridden with this message's ID.
func (m *Message) Reply(ctx context.Context, content string, opts ...MessageOptions) (*Message, error) {
	o := firstOpt(opts)
	o.ReplyToMessageIDs = []string{m.ID}
	return m.Channel.Send(ctx, content, o)
}

// EditOptions are the optional extras for editing a message.
type EditOptions struct {
	// HTMLEmbed replaces the message's HTML embed.
	HTMLEmbed string
	// Buttons replaces the message's buttons.
	Buttons []ButtonOption
}

// Edit replaces the message's content and returns the updated message. Only the
// bot's own messages can be edited.
func (m *Message) Edit(ctx context.Context, content string, opts ...EditOptions) (*Message, error) {
	var o EditOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	return m.client.editMessage(ctx, m.ChannelID, m.ID, content, o)
}

// Delete deletes the message.
func (m *Message) Delete(ctx context.Context) error {
	return m.client.deleteMessage(ctx, m.ChannelID, m.ID)
}

// Quote returns the raw token that quotes this message, e.g. "[q:837...]".
func (m *Message) Quote() string {
	return "[q:" + m.ID + "]"
}

// update applies a partial update (from a MessageUpdate event) in place.
func (m *Message) update(raw rawMessage) {
	if raw.Content != "" {
		m.Content = raw.Content
		m.Mentions = parseMentions(m.Content)
	}
	if raw.EditedAt != 0 {
		m.EditedAt = raw.EditedAt
	}
}
