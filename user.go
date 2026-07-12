package nerimity

// User is a global Nerimity user: the same User object is shared everywhere the
// user appears, independent of any server. For the per-server view (nickname,
// roles, permissions) see ServerMember.
type User struct {
	client *Client

	// ID is the user's unique snowflake ID.
	ID string
	// Username is the user's display name.
	Username string
	// Tag is the four-character discriminator shown after the username.
	Tag string
	// HexColor is the user's name colour, e.g. "#ffaa00".
	HexColor string
	// Badges is the user's badge bitfield.
	Badges int
	// Avatar is the avatar file path, empty if unset.
	Avatar string
	// Banner is the banner file path, empty if unset.
	Banner string
	// Bot reports whether the user is a bot account.
	Bot bool
	// JoinedAt is the account creation time as a Unix millisecond timestamp,
	// 0 if not supplied.
	JoinedAt int64
}

// Mention returns the raw token that mentions this user in message content,
// e.g. "[@:8374...]". Include it in a message you send to ping the user.
func (u *User) Mention() string {
	return "[@:" + u.ID + "]"
}

func newUser(client *Client, raw rawUser) *User {
	return &User{
		client:   client,
		ID:       raw.ID,
		Username: raw.Username,
		Tag:      raw.Tag,
		HexColor: raw.HexColor,
		Badges:   raw.Badges,
		Avatar:   raw.Avatar,
		Banner:   raw.Banner,
		Bot:      raw.Bot,
		JoinedAt: raw.JoinedAt,
	}
}

// ClientUser is the bot's own user, exposed as Client.User after the Ready
// event fires. It embeds User and adds self-only actions.
type ClientUser struct {
	*User
}
