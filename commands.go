package nerimity

import (
	"context"
	"fmt"
	"net/http"
)

// maxCommandDescriptionLength is the server-enforced limit on a command's
// description. It is checked client-side so you get a clear error before the
// request is made.
const maxCommandDescriptionLength = 60

// CommandDefinition describes a bot slash command to register with
// Client.UpdateCommands.
type CommandDefinition struct {
	// Name is the command name (used as "/name:botId ...").
	Name string `json:"name"`
	// Description is shown in the command list. Must be at most 60 characters.
	Description string `json:"description"`
	// Args is a human-readable hint for the command's arguments, e.g.
	// "<page number>".
	Args string `json:"args"`
	// Permissions, if non-zero, restricts who can use the command
	// (a RolePermission bitfield).
	Permissions int `json:"permissions,omitempty"`
}

// UpdateCommands registers the bot's slash commands, replacing any existing
// set. It authenticates with the bot token directly (not the logged-in
// session), so it can be called without an active connection.
//
// Each command description is limited to 60 characters; UpdateCommands returns
// an error naming the offending command before making any request if that limit
// is exceeded.
//
// This overwrites the full command list on every call, so pass every command
// you want to keep. You typically call this once when commands change, not on
// every startup.
func (c *Client) UpdateCommands(ctx context.Context, token string, commands []CommandDefinition) error {
	for _, cmd := range commands {
		if len(cmd.Description) > maxCommandDescriptionLength {
			return fmt.Errorf(
				"nerimity: command %q description is %d characters, exceeds the %d-character limit",
				cmd.Name, len(cmd.Description), maxCommandDescriptionLength,
			)
		}
	}
	url := c.apiBase + "/applications/bot/commands"
	body := map[string]any{"commands": commands}
	return c.doJSON(ctx, http.MethodPost, url, body, nil, token)
}
