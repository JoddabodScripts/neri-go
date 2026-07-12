package nerimity

import (
	"reflect"
	"testing"
)

func TestParseCommand(t *testing.T) {
	const botID = "999888777"

	tests := []struct {
		name    string
		content string
		botID   string
		want    *Command
	}{
		{
			name:    "matching command with args",
			content: "/help:999888777 page 2",
			botID:   botID,
			want:    &Command{Name: "help", Args: []string{"page", "2"}},
		},
		{
			name:    "matching command with no args",
			content: "/ping:999888777",
			botID:   botID,
			want:    &Command{Name: "ping", Args: []string{}},
		},
		{
			name:    "different bot id does not match",
			content: "/help:111111111 page 2",
			botID:   botID,
			want:    nil,
		},
		{
			name:    "not a command",
			content: "hello there",
			botID:   botID,
			want:    nil,
		},
		{
			name:    "empty bot id never matches",
			content: "/help: page 2",
			botID:   "",
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCommand(tt.content, tt.botID)
			if tt.want == nil {
				if got != nil {
					t.Errorf("parseCommand(%q) = %#v, want nil", tt.content, got)
				}
				return
			}
			if got == nil {
				t.Fatalf("parseCommand(%q) = nil, want %#v", tt.content, tt.want)
			}
			if got.Name != tt.want.Name || !reflect.DeepEqual(got.Args, tt.want.Args) {
				t.Errorf("parseCommand(%q) = %#v, want %#v", tt.content, got, tt.want)
			}
		})
	}
}
