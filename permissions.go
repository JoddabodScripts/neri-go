package nerimity

// RolePermission is a single Nerimity role permission bit. Permissions are
// stored as a bitfield on roles; combine and test them with the bitwise
// helpers HasBit, AddBit and RemoveBit.
type RolePermission int

// Role permission bits. These values match the Nerimity server (and the
// official JavaScript SDK) exactly. Do not renumber them.
const (
	// PermAdmin grants every permission and bypasses all other checks.
	PermAdmin RolePermission = 1
	// PermSendMessage allows sending messages in channels.
	PermSendMessage RolePermission = 2
	// PermManageRoles allows creating, editing and deleting roles.
	PermManageRoles RolePermission = 4
	// PermManageChannels allows creating, editing and deleting channels.
	PermManageChannels RolePermission = 8
	// PermKick allows kicking members from the server.
	PermKick RolePermission = 16
	// PermBan allows banning members from the server.
	PermBan RolePermission = 32
	// PermMentionEveryone allows using the @everyone mention.
	PermMentionEveryone RolePermission = 64
	// PermNicknameMember allows changing other members' nicknames.
	PermNicknameMember RolePermission = 128
	// PermMentionRoles allows mentioning roles.
	PermMentionRoles RolePermission = 256
)

// HasBit reports whether permissions contains bit.
func HasBit(permissions, bit int) bool {
	return permissions&bit == bit
}

// AddBit returns permissions with bit set.
func AddBit(permissions, bit int) int {
	return permissions | bit
}

// RemoveBit returns permissions with bit cleared.
func RemoveBit(permissions, bit int) int {
	return permissions &^ bit
}
