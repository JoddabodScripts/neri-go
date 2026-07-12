package nerimity

import (
	"testing"
	"time"
)

func TestParseEvent(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantName    string
		wantPayload string
		wantOK      bool
	}{
		{
			name:        "event with payload object",
			body:        `["message:created",{"message":{"id":"1"}}]`,
			wantName:    "message:created",
			wantPayload: `{"message":{"id":"1"}}`,
			wantOK:      true,
		},
		{
			name:     "event with no payload",
			body:     `["ping"]`,
			wantName: "ping",
			wantOK:   true,
		},
		{
			name:     "event with namespace prefix",
			body:     `/,["ready"]`,
			wantName: "ready",
			wantOK:   true,
		},
		{
			name:     "event with ack id prefix",
			body:     `12["ready"]`,
			wantName: "ready",
			wantOK:   true,
		},
		{
			name:   "malformed body",
			body:   `not json`,
			wantOK: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, payload, ok := parseEvent(tt.body)
			if ok != tt.wantOK {
				t.Fatalf("parseEvent(%q) ok = %v, want %v", tt.body, ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if name != tt.wantName {
				t.Errorf("parseEvent(%q) name = %q, want %q", tt.body, name, tt.wantName)
			}
			if tt.wantPayload != "" && string(payload) != tt.wantPayload {
				t.Errorf("parseEvent(%q) payload = %q, want %q", tt.body, payload, tt.wantPayload)
			}
		})
	}
}

func TestGatewayURL(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"https://nerimity.com", "wss://nerimity.com/socket.io/?EIO=4&transport=websocket"},
		{"http://localhost:8080", "ws://localhost:8080/socket.io/?EIO=4&transport=websocket"},
		{"https://nerimity.com/", "wss://nerimity.com/socket.io/?EIO=4&transport=websocket"},
	}
	for _, tt := range tests {
		got := gatewayURL(tt.in)
		if got != tt.want {
			t.Errorf("gatewayURL(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestBackoffDelayCapped(t *testing.T) {
	// After enough attempts, delay should be capped near 5s (allowing for
	// jitter up to +/-50%).
	d := backoffDelay(20)
	if d > 8*time.Second || d < 0 {
		t.Errorf("backoffDelay(20) = %v, want roughly within [0, 7.5s]", d)
	}
}
