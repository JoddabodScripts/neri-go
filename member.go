package nerimity

import "context"

// ServerMember is a user's per-server membership: their nickname, roles and
// computed permissions within one server. The global identity lives on the
// embedded-by-reference User (see the User field).
type ServerMember struct {
	client *Client

	// ID is the member's user ID (same as User.ID).
	ID string
	// User is the global user this membership belongs to.
	User *User
	// ServerID is the ID of the server.
	ServerID string
	// RoleIDs are the IDs of the roles assigned to this member.
	RoleIDs []string
	// Nickname is the member's server-specific nickname, empty if none.
	Nickname string
}

func newServerMember(client *Client, raw rawServerMember) *ServerMember {
	user, _ := client.users.get(raw.User.ID)
	return &ServerMember{
		client:   client,
		ID:       raw.User.ID,
		User:     user,
		ServerID: raw.ServerID,
		RoleIDs:  raw.RoleIDs,
		Nickname: raw.Nickname,
	}
}

// Server returns the server this membership belongs to, or nil if it is not
// cached.
func (m *ServerMember) Server() *Server {
	s, _ := m.client.servers.get(m.ServerID)
	return s
}

// Roles returns the member's roles, resolved from cache and skipping any that
// are not cached.
func (m *ServerMember) Roles() []*ServerRole {
	s := m.Server()
	if s == nil {
		return nil
	}
	roles := make([]*ServerRole, 0, len(m.RoleIDs))
	for _, id := range m.RoleIDs {
		if r, ok := s.roles.get(id); ok {
			roles = append(roles, r)
		}
	}
	return roles
}

// Permissions returns the member's effective permission bitfield: the server's
// default (@everyone) role permissions combined with every assigned role.
func (m *ServerMember) Permissions() int {
	s := m.Server()
	if s == nil {
		return 0
	}
	perms := 0
	if def, ok := s.roles.get(s.DefaultRoleID); ok {
		perms = def.Permissions
	}
	for _, r := range m.Roles() {
		perms = AddBit(perms, r.Permissions)
	}
	return perms
}

// HasPermission reports whether the member is allowed to perform the action
// guarded by perm. By default the server owner and anyone with PermAdmin are
// granted every permission; pass ignoreAdmin or ignoreCreator to bypass those
// shortcuts and test the raw bit only.
func (m *ServerMember) HasPermission(perm RolePermission, ignoreAdmin, ignoreCreator bool) bool {
	s := m.Server()
	if s == nil {
		return false
	}
	if !ignoreCreator && s.CreatedByID == m.ID {
		return true
	}
	perms := m.Permissions()
	if !ignoreAdmin && HasBit(perms, int(PermAdmin)) {
		return true
	}
	return HasBit(perms, int(perm))
}

// Kick removes the member from the server. Requires the PermKick permission.
func (m *ServerMember) Kick(ctx context.Context) error {
	return m.client.kickMember(ctx, m.ServerID, m.ID)
}

// Ban bans the member from the server. Requires the PermBan permission. reason
// may be empty.
func (m *ServerMember) Ban(ctx context.Context, reason string) error {
	return m.client.banMember(ctx, m.ServerID, m.ID, reason)
}

// Unban lifts a ban on this member's user. Requires the PermBan permission.
func (m *ServerMember) Unban(ctx context.Context) error {
	return m.client.unbanMember(ctx, m.ServerID, m.ID)
}

// Mention returns the raw token that mentions this member, e.g. "[@:837...]".
func (m *ServerMember) Mention() string {
	return "[@:" + m.ID + "]"
}
