package nerimity

import "context"

// Server is a Nerimity server (a "guild" in other platforms' terms). It owns
// collections of channels, members and roles.
type Server struct {
	client *Client

	// ID is the server's unique ID.
	ID string
	// Name is the server's display name.
	Name string
	// Avatar is the server avatar file path, empty if unset.
	Avatar string
	// DefaultRoleID is the ID of the @everyone role.
	DefaultRoleID string
	// CreatedByID is the user ID of the server owner.
	CreatedByID string

	channels *cache[*ServerChannel]
	members  *cache[*ServerMember]
	roles    *cache[*ServerRole]
}

func newServer(client *Client, raw rawServer) *Server {
	return &Server{
		client:        client,
		ID:            raw.ID,
		Name:          raw.Name,
		Avatar:        raw.Avatar,
		DefaultRoleID: raw.DefaultRoleID,
		CreatedByID:   raw.CreatedByID,
		channels:      newCache[*ServerChannel](0),
		members:       newCache[*ServerMember](0),
		roles:         newCache[*ServerRole](0),
	}
}

// Channels returns a snapshot of the server's cached channels.
func (s *Server) Channels() []*ServerChannel { return s.channels.values() }

// Members returns a snapshot of the server's cached members.
func (s *Server) Members() []*ServerMember { return s.members.values() }

// Roles returns a snapshot of the server's cached roles.
func (s *Server) Roles() []*ServerRole { return s.roles.values() }

// Member returns the cached member for the given user ID, or nil.
func (s *Server) Member(userID string) *ServerMember {
	m, _ := s.members.get(userID)
	return m
}

// Role returns the cached role for the given role ID, or nil.
func (s *Server) Role(roleID string) *ServerRole {
	r, _ := s.roles.get(roleID)
	return r
}

// Channel returns the cached channel for the given channel ID, or nil.
func (s *Server) Channel(channelID string) *ServerChannel {
	c, _ := s.channels.get(channelID)
	return c
}

// KickMember kicks a user from the server by ID. Requires PermKick.
func (s *Server) KickMember(ctx context.Context, userID string) error {
	return s.client.kickMember(ctx, s.ID, userID)
}

// BanMember bans a user from the server by ID. Requires PermBan. reason may be
// empty.
func (s *Server) BanMember(ctx context.Context, userID, reason string) error {
	return s.client.banMember(ctx, s.ID, userID, reason)
}

// UnbanMember lifts a ban by user ID. Requires PermBan.
func (s *Server) UnbanMember(ctx context.Context, userID string) error {
	return s.client.unbanMember(ctx, s.ID, userID)
}
