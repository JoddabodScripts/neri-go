package nerimity

import (
	"context"
	"sync"
)

// ButtonOption describes a button to attach to a message.
type ButtonOption struct {
	// ID is the identifier reported back in the MessageButtonClick event.
	ID string `json:"id"`
	// Label is the text shown on the button.
	Label string `json:"label"`
	// Alert, when true, makes the button show a confirmation alert before the
	// click is delivered.
	Alert bool `json:"alert,omitempty"`
}

// MessageOptions are the optional extras for sending a message. The zero value
// is valid: sending with no options at all is the common case.
type MessageOptions struct {
	// HTMLEmbed is an HTML embed rendered below the message. Any user-supplied
	// text inside it MUST be passed through EscapeHTML first; see html.go.
	HTMLEmbed string
	// NerimityCdnFileID attaches a previously uploaded CDN file (see
	// AttachmentBuilder).
	NerimityCdnFileID string
	// Buttons attaches interactive buttons to the message.
	Buttons []ButtonOption
	// Silent suppresses notifications for this message.
	Silent bool
	// ReplyToMessageIDs marks this message as a reply to the given messages.
	ReplyToMessageIDs []string
	// MentionReplies pings the authors of the replied-to messages.
	MentionReplies bool
}

func firstOpt(opts []MessageOptions) MessageOptions {
	if len(opts) > 0 {
		return opts[0]
	}
	return MessageOptions{}
}

// Channel is a channel messages can be sent to. A channel with a non-empty
// ServerID is a server channel; otherwise it is a DM. ServerChannel is an alias
// for this type used to document server-channel event handlers.
type Channel struct {
	client *Client

	// ID is the channel's unique ID.
	ID string
	// Type is the channel kind (DM, server text or category).
	Type ChannelType
	// ServerID is the owning server's ID, empty for DM channels.
	ServerID string
	// Name is the channel name, empty for DM channels.
	Name string
	// Permissions is the channel permission bitfield (server channels only).
	Permissions int
	// CategoryID is the parent category's ID, if any.
	CategoryID string
	// Order is the channel's position within its category or server.
	Order int
	// CreatedByID is the user ID of the channel's creator, if known.
	CreatedByID string
	// CreatedAt is the creation time as a Unix millisecond timestamp.
	CreatedAt int64
	// LastMessagedAt is the last-activity time as a Unix millisecond timestamp.
	LastMessagedAt int64

	// sendMu serialises sends on this channel so that message order is
	// preserved and requests never overlap, matching the JavaScript SDK's
	// per-channel AsyncFunctionQueue.
	sendMu sync.Mutex
}

// ServerChannel is a Channel that belongs to a server. It is an alias of
// Channel: every field the JavaScript SDK's ServerChannel exposes (Name,
// ServerID, Permissions, CategoryID) already lives on Channel. The alias keeps
// server-channel event handler signatures self-documenting.
type ServerChannel = Channel

func newChannel(client *Client, raw rawChannel) *Channel {
	return &Channel{
		client:         client,
		ID:             raw.ID,
		Type:           raw.Type,
		ServerID:       raw.ServerID,
		Name:           raw.Name,
		Permissions:    raw.Permissions,
		CategoryID:     raw.CategoryID,
		Order:          raw.Order,
		CreatedByID:    raw.CreatedByID,
		CreatedAt:      raw.CreatedAt,
		LastMessagedAt: raw.LastMessage,
	}
}

// Server returns the server this channel belongs to, or nil for a DM channel or
// when the server is not cached.
func (c *Channel) Server() *Server {
	if c.ServerID == "" {
		return nil
	}
	s, _ := c.client.servers.get(c.ServerID)
	return s
}

// Send posts a message to the channel and returns the created message. Sends on
// the same channel are serialised, preserving order. opts is optional; pass at
// most one MessageOptions.
func (c *Channel) Send(ctx context.Context, content string, opts ...MessageOptions) (*Message, error) {
	o := firstOpt(opts)
	c.sendMu.Lock()
	defer c.sendMu.Unlock()
	return c.client.postMessage(ctx, c.ID, content, o)
}

// DeleteMessage deletes a message in this channel by ID.
func (c *Channel) DeleteMessage(ctx context.Context, messageID string) error {
	return c.client.deleteMessage(ctx, c.ID, messageID)
}

// Mention returns the raw token that links to this channel, e.g. "[#:837...]".
func (c *Channel) Mention() string {
	return "[#:" + c.ID + "]"
}
