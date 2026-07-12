package nerimity

import "strings"

// EscapeHTML escapes the five HTML-significant characters (& < > " ') so that a
// user-supplied string can be safely interpolated into an HTMLEmbed without
// breaking the markup.
//
// This is not cosmetic. Nerimity's server-side HTML embed validator counts
// opening and closing tags with a regular expression to check that the markup
// is balanced. A raw "<" or ">" inside otherwise plain text is seen by that
// regex as a tag and throws the count off, causing the whole embed to be
// rejected. Any value you did not write yourself (usernames, message content,
// API responses) MUST be passed through EscapeHTML before it goes into an
// embed. See docs/messages.md for the full list of validator gotchas.
func EscapeHTML(s string) string {
	return htmlEscaper.Replace(s)
}

// Order matters: & must be replaced first so we don't double-escape the
// ampersands introduced by the other replacements.
var htmlEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&#39;",
)
