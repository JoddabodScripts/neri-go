package nerimity

// This file mirrors the wire structs from the Nerimity JavaScript SDK's
// RawData.ts. These types are unexported: they exist only to unmarshal
// WebSocket and REST payloads, which are then converted into the friendly
// exported types (User, Message, Server, ...). If the server protocol changes,
// this is the file that must be re-verified against ~/nrepos/nerimity.js.

// MessageType identifies why a message exists (normal content vs. a system
// event like a member joining).
type MessageType int

const (
	MessageTypeContent     MessageType = 0
	MessageTypeJoinServer  MessageType = 1
	MessageTypeLeaveServer MessageType = 2
	MessageTypeKickUser    MessageType = 3
	MessageTypeBanUser     MessageType = 4
)

// ChannelType identifies whether a channel is a DM, a server text channel or a
// category.
type ChannelType int

const (
	ChannelTypeDMText     ChannelType = 0
	ChannelTypeServerText ChannelType = 1
	ChannelTypeCategory   ChannelType = 2
)

type rawUser struct {
	ID        string `json:"id"`
	Avatar    string `json:"avatar,omitempty"`
	Banner    string `json:"banner,omitempty"`
	Username  string `json:"username"`
	Bot       bool   `json:"bot,omitempty"`
	HexColor  string `json:"hexColor"`
	Tag       string `json:"tag"`
	Badges    int    `json:"badges"`
	JoinedAt  int64  `json:"joinedAt,omitempty"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

type rawServer struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	HexColor         string `json:"hexColor"`
	DefaultChannelID string `json:"defaultChannelId"`
	SystemChannelID  string `json:"systemChannelId,omitempty"`
	Avatar           string `json:"avatar,omitempty"`
	Banner           string `json:"banner,omitempty"`
	DefaultRoleID    string `json:"defaultRoleId"`
	CreatedByID      string `json:"createdById"`
	CreatedAt        int64  `json:"createdAt"`
	Verified         bool   `json:"verified"`
}

type rawServerMember struct {
	ServerID string   `json:"serverId"`
	User     rawUser  `json:"user"`
	Nickname string   `json:"nickname,omitempty"`
	JoinedAt int64    `json:"joinedAt"`
	RoleIDs  []string `json:"roleIds"`
}

type rawServerRole struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon,omitempty"`
	Order       int    `json:"order"`
	HexColor    string `json:"hexColor"`
	CreatedByID string `json:"createdById"`
	Permissions int    `json:"permissions"`
	ServerID    string `json:"serverId"`
	HideRole    bool   `json:"hideRole"`
	BotRole     bool   `json:"botRole,omitempty"`
	ApplyOnJoin bool   `json:"applyOnJoin,omitempty"`
}

type rawChannel struct {
	ID          string      `json:"id"`
	CategoryID  string      `json:"categoryId,omitempty"`
	Name        string      `json:"name"`
	CreatedByID string      `json:"createdById,omitempty"`
	ServerID    string      `json:"serverId,omitempty"`
	Type        ChannelType `json:"type"`
	Permissions int         `json:"permissions,omitempty"`
	CreatedAt   int64       `json:"createdAt"`
	LastMessage int64       `json:"lastMessagedAt,omitempty"`
	Order       int         `json:"order,omitempty"`
}

type rawMessageButton struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Alert bool   `json:"alert,omitempty"`
}

type rawMessageReaction struct {
	Name    string `json:"name"`
	EmojiID string `json:"emojiId,omitempty"`
	Gif     bool   `json:"gif,omitempty"`
	Reacted bool   `json:"reacted"`
	Count   int    `json:"count"`
}

type rawReplyMessage struct {
	ReplyToMessage *rawMessage `json:"replyToMessage,omitempty"`
}

type rawMessage struct {
	ID             string               `json:"id"`
	ChannelID      string               `json:"channelId"`
	Silent         bool                 `json:"silent,omitempty"`
	Content        string               `json:"content,omitempty"`
	CreatedBy      rawUser              `json:"createdBy"`
	Type           MessageType          `json:"type"`
	CreatedAt      int64                `json:"createdAt"`
	Pinned         bool                 `json:"pinned,omitempty"`
	EditedAt       int64                `json:"editedAt,omitempty"`
	Mentions       []rawUser            `json:"mentions,omitempty"`
	Attachments    []rawAttachment      `json:"attachments,omitempty"`
	Reactions      []rawMessageReaction `json:"reactions,omitempty"`
	HTMLEmbed      string               `json:"htmlEmbed,omitempty"`
	MentionReplies bool                 `json:"mentionReplies,omitempty"`
	ReplyMessages  []rawReplyMessage    `json:"replyMessages,omitempty"`
	Buttons        []rawMessageButton   `json:"buttons,omitempty"`
	WebhookID      string               `json:"webhookId,omitempty"`
}

type rawAttachment struct {
	ID       string `json:"id"`
	Provider string `json:"provider,omitempty"`
	FileID   string `json:"fileId,omitempty"`
	Mime     string `json:"mime,omitempty"`
	Path     string `json:"path,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	Filesize int64  `json:"filesize,omitempty"`
}

type rawCDNUpload struct {
	FileID string `json:"fileId"`
}

// ---- WebSocket event payloads ----

type authenticatedPayload struct {
	User          rawUser           `json:"user"`
	Servers       []rawServer       `json:"servers"`
	ServerMembers []rawServerMember `json:"serverMembers"`
	Channels      []rawChannel      `json:"channels"`
	ServerRoles   []rawServerRole   `json:"serverRoles"`
}

type messageButtonClickPayload struct {
	MessageID string            `json:"messageId"`
	ChannelID string            `json:"channelId"`
	ButtonID  string            `json:"buttonId"`
	UserID    string            `json:"userId"`
	Type      string            `json:"type"`
	Data      map[string]string `json:"data,omitempty"`
}

type reactionAddedPayload struct {
	MessageID       string `json:"messageId"`
	ChannelID       string `json:"channelId"`
	Count           int    `json:"count"`
	ReactedByUserID string `json:"reactedByUserId"`
	EmojiID         string `json:"emojiId,omitempty"`
	Name            string `json:"name"`
	Gif             bool   `json:"gif,omitempty"`
}

type reactionRemovedPayload struct {
	MessageID             string `json:"messageId"`
	ChannelID             string `json:"channelId"`
	Count                 int    `json:"count"`
	ReactionRemovedByUser string `json:"reactionRemovedByUserId"`
	EmojiID               string `json:"emojiId,omitempty"`
	Name                  string `json:"name"`
	Gif                   bool   `json:"gif,omitempty"`
}
