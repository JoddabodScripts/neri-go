package nerimity

import (
	"regexp"
	"strings"
)

// MentionType classifies a mention token found in message content.
type MentionType string

const (
	// MentionTypeUser is a mention of a specific user: "[@:userId]".
	MentionTypeUser MentionType = "user"
	// MentionTypeEveryone is the @everyone mention: "[@:e]".
	MentionTypeEveryone MentionType = "everyone"
	// MentionTypeSomeone is the @someone mention: "[@:s]".
	MentionTypeSomeone MentionType = "someone"
)

// Mention is a single mention token parsed from message content. Following the
// JavaScript SDK, tokens are reported raw and are not resolved to display
// names; for user mentions, look up UserID in the client's user cache if you
// need the User object.
type Mention struct {
	// Type is the kind of mention.
	Type MentionType
	// UserID is the mentioned user's ID, set only when Type is
	// MentionTypeUser.
	UserID string
	// Raw is the exact token as it appeared in the content, e.g. "[@:837...]".
	Raw string
}

// mentionRegex matches a user mention ("[@:123]"), everyone ("[@:e]") or
// someone ("[@:s]").
var mentionRegex = regexp.MustCompile(`\[@:(e|s|[0-9]+)\]`)

// parseMentions extracts every mention token from content in order of
// appearance.
func parseMentions(content string) []Mention {
	matches := mentionRegex.FindAllStringSubmatch(content, -1)
	if matches == nil {
		return nil
	}
	out := make([]Mention, 0, len(matches))
	for _, m := range matches {
		mention := Mention{Raw: m[0]}
		switch m[1] {
		case "e":
			mention.Type = MentionTypeEveryone
		case "s":
			mention.Type = MentionTypeSomeone
		default:
			mention.Type = MentionTypeUser
			mention.UserID = m[1]
		}
		out = append(out, mention)
	}
	return out
}

// StripMentions removes every mention token ("[@:...]") from content and
// collapses the whitespace left behind. Use it when a bot wants to process the
// plain text of a message without mention markup.
func StripMentions(content string) string {
	stripped := mentionRegex.ReplaceAllString(content, "")
	return strings.TrimSpace(collapseSpaces.ReplaceAllString(stripped, " "))
}

var collapseSpaces = regexp.MustCompile(`[ \t]{2,}`)
