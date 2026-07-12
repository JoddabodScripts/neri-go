# Permissions

Server roles carry a permission bitfield. These bit values are fixed by the
Nerimity server and must match exactly:

| Constant              | Value | Grants                          |
|------------------------|-------|----------------------------------|
| `PermAdmin`            | 1     | Everything, bypassing all other checks |
| `PermSendMessage`      | 2     | Sending messages                |
| `PermManageRoles`      | 4     | Creating/editing/deleting roles |
| `PermManageChannels`   | 8     | Creating/editing/deleting channels |
| `PermKick`             | 16    | Kicking members                 |
| `PermBan`              | 32    | Banning members                 |
| `PermMentionEveryone`  | 64    | Using `[@:e]`                   |
| `PermNicknameMember`   | 128   | Changing other members' nicknames |
| `PermMentionRoles`     | 256   | Mentioning roles                |

## Bitwise helpers

```go
nerimity.HasBit(permissions, bit int) bool
nerimity.AddBit(permissions, bit int) int
nerimity.RemoveBit(permissions, bit int) int
```

```go
perms := 0
perms = nerimity.AddBit(perms, int(nerimity.PermSendMessage))
perms = nerimity.AddBit(perms, int(nerimity.PermKick))

nerimity.HasBit(perms, int(nerimity.PermKick)) // true
```

## Checking a role directly

```go
if role.HasPermission(nerimity.PermManageChannels) {
	// ...
}
```

This checks the role's own bitfield only — it does not account for
admin/owner overrides.

## Checking what a member can actually do

Use `ServerMember.HasPermission`, which combines the server's default
(`@everyone`) role with every role the member has, and short-circuits for the
server owner and admins:

```go
member := server.Member(userID)

if member.HasPermission(nerimity.PermKick, false, false) {
	// member can kick: they're the owner, an admin, or have PermKick directly
}
```

The second and third arguments (`ignoreAdmin`, `ignoreCreator`) let you bypass
those shortcuts when you specifically need to test the raw bit, e.g. to check
"does this role literally have Kick" regardless of who holds it:

```go
member.HasPermission(nerimity.PermKick, true, true) // no owner/admin shortcut
```

`ServerMember.Permissions()` returns the raw combined bitfield if you want to
do your own bit tests instead.
