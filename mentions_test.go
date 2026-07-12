package nerimity

import (
	"reflect"
	"testing"
)

func TestParseMentions(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []Mention
	}{
		{
			name:    "no mentions",
			content: "hello world",
			want:    nil,
		},
		{
			name:    "single user mention",
			content: "hey [@:12345] how's it going",
			want: []Mention{
				{Type: MentionTypeUser, UserID: "12345", Raw: "[@:12345]"},
			},
		},
		{
			name:    "everyone and someone",
			content: "[@:e] and [@:s] listen up",
			want: []Mention{
				{Type: MentionTypeEveryone, Raw: "[@:e]"},
				{Type: MentionTypeSomeone, Raw: "[@:s]"},
			},
		},
		{
			name:    "multiple user mentions in order",
			content: "[@:1] [@:2] [@:3]",
			want: []Mention{
				{Type: MentionTypeUser, UserID: "1", Raw: "[@:1]"},
				{Type: MentionTypeUser, UserID: "2", Raw: "[@:2]"},
				{Type: MentionTypeUser, UserID: "3", Raw: "[@:3]"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseMentions(tt.content)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseMentions(%q) = %#v, want %#v", tt.content, got, tt.want)
			}
		})
	}
}

func TestStripMentions(t *testing.T) {
	tests := []struct {
		content string
		want    string
	}{
		{"hey [@:12345] how are you", "hey how are you"},
		{"[@:e] listen up", "listen up"},
		{"no mentions here", "no mentions here"},
		{"[@:1][@:2] adjacent", "adjacent"},
	}
	for _, tt := range tests {
		got := StripMentions(tt.content)
		if got != tt.want {
			t.Errorf("StripMentions(%q) = %q, want %q", tt.content, got, tt.want)
		}
	}
}
