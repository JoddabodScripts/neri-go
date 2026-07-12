package nerimity

import "testing"

func TestHasBit(t *testing.T) {
	perms := int(PermSendMessage) | int(PermManageRoles)
	if !HasBit(perms, int(PermSendMessage)) {
		t.Error("expected PermSendMessage to be set")
	}
	if !HasBit(perms, int(PermManageRoles)) {
		t.Error("expected PermManageRoles to be set")
	}
	if HasBit(perms, int(PermAdmin)) {
		t.Error("did not expect PermAdmin to be set")
	}
	if HasBit(perms, int(PermBan)) {
		t.Error("did not expect PermBan to be set")
	}
}

func TestAddBit(t *testing.T) {
	perms := AddBit(0, int(PermSendMessage))
	perms = AddBit(perms, int(PermKick))
	if !HasBit(perms, int(PermSendMessage)) || !HasBit(perms, int(PermKick)) {
		t.Errorf("AddBit did not set both bits, got %d", perms)
	}
}

func TestRemoveBit(t *testing.T) {
	perms := int(PermAdmin) | int(PermSendMessage) | int(PermBan)
	perms = RemoveBit(perms, int(PermSendMessage))
	if HasBit(perms, int(PermSendMessage)) {
		t.Error("expected PermSendMessage to be cleared")
	}
	if !HasBit(perms, int(PermAdmin)) || !HasBit(perms, int(PermBan)) {
		t.Error("RemoveBit cleared unrelated bits")
	}
}

func TestPermissionValues(t *testing.T) {
	// These values must match the Nerimity server and the JS SDK exactly.
	cases := map[RolePermission]int{
		PermAdmin:           1,
		PermSendMessage:     2,
		PermManageRoles:     4,
		PermManageChannels:  8,
		PermKick:            16,
		PermBan:             32,
		PermMentionEveryone: 64,
		PermNicknameMember:  128,
		PermMentionRoles:    256,
	}
	for perm, want := range cases {
		if int(perm) != want {
			t.Errorf("permission %v = %d, want %d", perm, perm, want)
		}
	}
}
