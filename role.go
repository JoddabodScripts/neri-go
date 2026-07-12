package nerimity

// ServerRole is a role within a server. Roles carry a permission bitfield;
// test it with HasBit and the Perm* constants.
type ServerRole struct {
	client *Client

	// ID is the role's unique ID.
	ID string
	// Name is the role's display name.
	Name string
	// Permissions is the role's permission bitfield.
	Permissions int
	// HexColor is the role's colour, e.g. "#ffaa00".
	HexColor string
	// Order is the role's position in the hierarchy.
	Order int
	// ServerID is the ID of the server this role belongs to.
	ServerID string
	// IsDefaultRole reports whether this is the server's @everyone role.
	IsDefaultRole bool
}

func newServerRole(client *Client, raw rawServerRole) *ServerRole {
	r := &ServerRole{
		client:      client,
		ID:          raw.ID,
		Name:        raw.Name,
		Permissions: raw.Permissions,
		HexColor:    raw.HexColor,
		Order:       raw.Order,
		ServerID:    raw.ServerID,
	}
	if s, ok := client.servers.get(raw.ServerID); ok {
		r.IsDefaultRole = s.DefaultRoleID == r.ID
	}
	return r
}

// HasPermission reports whether the role's own bitfield includes perm. Note
// this checks the role in isolation; to check what a member can actually do
// (accounting for @everyone, admin and server ownership) use
// ServerMember.HasPermission.
func (r *ServerRole) HasPermission(perm RolePermission) bool {
	return HasBit(r.Permissions, int(perm))
}
