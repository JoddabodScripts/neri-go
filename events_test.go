package nerimity

import (
	"encoding/json"
	"testing"
)

// TestOnAuthenticatedSurvivesMalformedSubField is a regression test for a bug
// where a single malformed field anywhere in the "user:authenticated" payload
// (for example, a serverRoles entry shaped differently than expected) caused
// json.Unmarshal to fail for the whole payload. Because that failure was
// handled with a silent early return, Client.user was never set, OnReady
// never fired, and every later event handler that called Client.User() panicked
// with a nil pointer dereference on the *ClientUser embedding.
//
// The "user" field must be decoded independently of the rest, so it survives
// even when other parts of the payload don't match this SDK's expectations.
func TestOnAuthenticatedSurvivesMalformedSubField(t *testing.T) {
	c := New(Options{})

	var gotErr error
	c.OnError(func(err error) { gotErr = err })

	var readyFired bool
	c.OnReady(func() { readyFired = true })

	// serverRoles is malformed: a string where an array of role objects is
	// expected. This must not prevent "user" from being decoded.
	payload := json.RawMessage(`{
		"user": {"id": "1", "username": "bot", "tag": "0001", "hexColor": "#fff", "badges": 0},
		"servers": [],
		"channels": [],
		"serverMembers": [],
		"serverRoles": "this is not an array"
	}`)

	c.onAuthenticated(payload)

	if c.User() == nil {
		t.Fatal("Client.User() is nil after onAuthenticated, want it set despite the malformed serverRoles field")
	}
	if c.User().ID != "1" {
		t.Errorf("Client.User().ID = %q, want %q", c.User().ID, "1")
	}
	if !readyFired {
		t.Error("OnReady handler did not fire despite a valid \"user\" field")
	}
	if gotErr == nil {
		t.Error("OnError handler did not fire for the malformed serverRoles field")
	}
}

// TestOnAuthenticatedReportsAndSkipsWhenUserFieldMissing verifies that when
// the "user" field itself is unparseable, the client reports the error and
// leaves Client.User() nil rather than panicking or silently pretending to be
// ready.
func TestOnAuthenticatedReportsAndSkipsWhenUserFieldMissing(t *testing.T) {
	c := New(Options{})

	var gotErr error
	c.OnError(func(err error) { gotErr = err })

	var readyFired bool
	c.OnReady(func() { readyFired = true })

	payload := json.RawMessage(`{"user": "not a user object"}`)
	c.onAuthenticated(payload)

	if c.User() != nil {
		t.Error("Client.User() should remain nil when the \"user\" field itself is unparseable")
	}
	if readyFired {
		t.Error("OnReady should not fire when the \"user\" field failed to decode")
	}
	if gotErr == nil {
		t.Error("OnError handler did not fire for the unparseable \"user\" field")
	}
}

// TestOnAuthenticatedSuccess is a sanity check that a well-formed payload
// still works end to end: user set, collections populated, Ready fires, no
// error reported.
func TestOnAuthenticatedSuccess(t *testing.T) {
	c := New(Options{})

	var gotErr error
	c.OnError(func(err error) { gotErr = err })

	var readyFired bool
	c.OnReady(func() { readyFired = true })

	payload := json.RawMessage(`{
		"user": {"id": "1", "username": "bot", "tag": "0001", "hexColor": "#fff", "badges": 0},
		"servers": [{"id": "s1", "name": "Test", "hexColor": "#000", "defaultChannelId": "c1", "defaultRoleId": "r1", "createdById": "1", "createdAt": 1, "verified": false}],
		"channels": [{"id": "c1", "name": "general", "serverId": "s1", "type": 1, "createdAt": 1}],
		"serverMembers": [],
		"serverRoles": [{"id": "r1", "name": "everyone", "order": 1, "hexColor": "#000", "createdById": "1", "permissions": 2, "serverId": "s1", "hideRole": false}]
	}`)

	c.onAuthenticated(payload)

	if c.User() == nil || c.User().ID != "1" {
		t.Fatalf("Client.User() = %v, want user with ID 1", c.User())
	}
	if !readyFired {
		t.Error("OnReady did not fire")
	}
	if gotErr != nil {
		t.Errorf("unexpected error reported: %v", gotErr)
	}
	server := c.Server("s1")
	if server == nil {
		t.Fatal("server s1 not cached")
	}
	if server.Channel("c1") == nil {
		t.Error("channel c1 not attached to server s1")
	}
	if server.Role("r1") == nil {
		t.Error("role r1 not cached on server s1")
	}
}
