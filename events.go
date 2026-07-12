package nerimity

import (
	"encoding/json"
	"fmt"
	"sync"
)

// eventHandlers holds every registered callback. Each event has its own slice
// so multiple handlers can be registered for the same event, matching the
// JavaScript SDK's EventEmitter semantics.
type eventHandlers struct {
	mu sync.RWMutex

	onReady                  []func()
	onMessageCreate          []func(*Message)
	onMessageUpdate          []func(*Message)
	onMessageDelete          []func(MessageDeleteEvent)
	onMessageReactionAdded   []func(*Reaction, *User)
	onMessageReactionRemoved []func(*Reaction, *User)
	onMessageButtonClick     []func(*MessageButton)
	onServerMemberJoined     []func(*ServerMember)
	onServerMemberLeft       []func(*ServerMember)
	onServerMemberUpdated    []func(*ServerMember)
	onServerJoined           []func(*Server)
	onServerLeft             []func(*Server)
	onServerChannelCreated   []func(*ServerChannel)
	onServerChannelUpdated   []func(*ServerChannel)
	onServerChannelDeleted   []func(ServerChannelDeleteEvent)
	onServerRoleCreated      []func(*ServerRole)
	onServerRoleUpdated      []func(*ServerRole)
	onServerRoleDeleted      []func(*ServerRole)
	onServerRoleOrderUpdated []func(server *Server, roles []*ServerRole)
	onError                  []func(error)
}

// MessageDeleteEvent is the payload for OnMessageDelete: the message is already
// gone from cache by the time the event fires, so only its IDs are available.
type MessageDeleteEvent struct {
	MessageID string
	ChannelID string
}

// ServerChannelDeleteEvent is the payload for OnServerChannelDeleted.
type ServerChannelDeleteEvent struct {
	ServerID  string
	ChannelID string
}

// OnReady registers a handler called once the client has authenticated and
// populated its caches.
func (c *Client) OnReady(fn func()) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onReady = append(c.handlers.onReady, fn)
}

// OnMessageCreate registers a handler called when a message is sent.
func (c *Client) OnMessageCreate(fn func(*Message)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onMessageCreate = append(c.handlers.onMessageCreate, fn)
}

// OnMessageUpdate registers a handler called when a message is edited.
func (c *Client) OnMessageUpdate(fn func(*Message)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onMessageUpdate = append(c.handlers.onMessageUpdate, fn)
}

// OnMessageDelete registers a handler called when a message is deleted.
func (c *Client) OnMessageDelete(fn func(MessageDeleteEvent)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onMessageDelete = append(c.handlers.onMessageDelete, fn)
}

// OnMessageReactionAdded registers a handler called when a reaction is added to
// a message. user is nil if the reacting user is not cached.
func (c *Client) OnMessageReactionAdded(fn func(*Reaction, *User)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onMessageReactionAdded = append(c.handlers.onMessageReactionAdded, fn)
}

// OnMessageReactionRemoved registers a handler called when a reaction is
// removed from a message. user is nil if the user is not cached.
func (c *Client) OnMessageReactionRemoved(fn func(*Reaction, *User)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onMessageReactionRemoved = append(c.handlers.onMessageReactionRemoved, fn)
}

// OnMessageButtonClick registers a handler called when a user clicks a button
// on a message. Call button.Respond to reply.
func (c *Client) OnMessageButtonClick(fn func(*MessageButton)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onMessageButtonClick = append(c.handlers.onMessageButtonClick, fn)
}

// OnServerMemberJoined registers a handler called when a member joins a server.
func (c *Client) OnServerMemberJoined(fn func(*ServerMember)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onServerMemberJoined = append(c.handlers.onServerMemberJoined, fn)
}

// OnServerMemberLeft registers a handler called when a member leaves a server.
func (c *Client) OnServerMemberLeft(fn func(*ServerMember)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onServerMemberLeft = append(c.handlers.onServerMemberLeft, fn)
}

// OnServerMemberUpdated registers a handler called when a member's roles
// change.
func (c *Client) OnServerMemberUpdated(fn func(*ServerMember)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onServerMemberUpdated = append(c.handlers.onServerMemberUpdated, fn)
}

// OnServerJoined registers a handler called when the bot joins a server.
func (c *Client) OnServerJoined(fn func(*Server)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onServerJoined = append(c.handlers.onServerJoined, fn)
}

// OnServerLeft registers a handler called when the bot leaves (or is removed
// from) a server.
func (c *Client) OnServerLeft(fn func(*Server)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onServerLeft = append(c.handlers.onServerLeft, fn)
}

// OnServerChannelCreated registers a handler called when a server channel is
// created.
func (c *Client) OnServerChannelCreated(fn func(*ServerChannel)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onServerChannelCreated = append(c.handlers.onServerChannelCreated, fn)
}

// OnServerChannelUpdated registers a handler called when a server channel is
// updated.
func (c *Client) OnServerChannelUpdated(fn func(*ServerChannel)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onServerChannelUpdated = append(c.handlers.onServerChannelUpdated, fn)
}

// OnServerChannelDeleted registers a handler called when a server channel is
// deleted.
func (c *Client) OnServerChannelDeleted(fn func(ServerChannelDeleteEvent)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onServerChannelDeleted = append(c.handlers.onServerChannelDeleted, fn)
}

// OnServerRoleCreated registers a handler called when a role is created.
func (c *Client) OnServerRoleCreated(fn func(*ServerRole)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onServerRoleCreated = append(c.handlers.onServerRoleCreated, fn)
}

// OnServerRoleUpdated registers a handler called when a role is updated.
func (c *Client) OnServerRoleUpdated(fn func(*ServerRole)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onServerRoleUpdated = append(c.handlers.onServerRoleUpdated, fn)
}

// OnServerRoleDeleted registers a handler called when a role is deleted.
func (c *Client) OnServerRoleDeleted(fn func(*ServerRole)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onServerRoleDeleted = append(c.handlers.onServerRoleDeleted, fn)
}

// OnServerRoleOrderUpdated registers a handler called when a server's role
// order changes. roles is every cached role for the server, in the new order.
func (c *Client) OnServerRoleOrderUpdated(fn func(server *Server, roles []*ServerRole)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onServerRoleOrderUpdated = append(c.handlers.onServerRoleOrderUpdated, fn)
}

// OnError registers a handler called whenever the client fails to decode a
// gateway event's payload. This should never happen in normal operation; if it
// does, the gateway sent something this SDK version doesn't understand for
// that event, and the event in question was dropped (its handlers did not
// run, and any cache updates it would have made did not happen). Registering
// a handler here is optional but strongly recommended in production, since a
// silently dropped "user:authenticated" payload, for example, otherwise looks
// like the bot connecting successfully but every later handler seeing a nil
// Client.User().
func (c *Client) OnError(fn func(error)) {
	c.handlers.mu.Lock()
	defer c.handlers.mu.Unlock()
	c.handlers.onError = append(c.handlers.onError, fn)
}

// reportError delivers err to every registered OnError handler. Safe to call
// with no handlers registered (a no-op).
func (c *Client) reportError(err error) {
	for _, fn := range c.snapshotHandlers().onError {
		fn(err)
	}
}

// ---- dispatch ----

// Gateway event names, verified against nerimity.js's SocketServerEvents.
const (
	wsAuthenticated         = "user:authenticated"
	wsServerMemberJoined    = "server:member_joined"
	wsServerMemberUpdated   = "server:member_updated"
	wsServerMemberLeft      = "server:member_left"
	wsServerJoined          = "server:joined"
	wsServerChannelCreated  = "server:channel_created"
	wsServerChannelUpdated  = "server:channel_updated"
	wsServerChannelDeleted  = "server:channel_deleted"
	wsServerLeft            = "server:left"
	wsMessageCreated        = "message:created"
	wsMessageUpdated        = "message:updated"
	wsMessageDeleted        = "message:deleted"
	wsMessageButtonClicked  = "message:button_clicked"
	wsServerRoleCreated     = "server:role_created"
	wsServerRoleDeleted     = "server:role_deleted"
	wsServerRoleUpdated     = "server:role_updated"
	wsServerRoleOrderUpdate = "server:role_order_updated"
	wsMessageReactionAdded  = "message:reaction_added"
	wsMessageReactionRemove = "message:reaction_removed"
)

func (c *Client) handleEvent(ev incomingEvent) {
	switch ev.name {
	case wsAuthenticated:
		c.onAuthenticated(ev.payload)
	case wsServerMemberJoined:
		c.onServerMemberJoined(ev.payload)
	case wsServerMemberUpdated:
		c.onServerMemberUpdated(ev.payload)
	case wsServerMemberLeft:
		c.onServerMemberLeft(ev.payload)
	case wsServerJoined:
		c.onServerJoined(ev.payload)
	case wsServerChannelCreated:
		c.onServerChannelCreated(ev.payload)
	case wsServerChannelUpdated:
		c.onServerChannelUpdated(ev.payload)
	case wsServerChannelDeleted:
		c.onServerChannelDeleted(ev.payload)
	case wsServerLeft:
		c.onServerLeft(ev.payload)
	case wsMessageCreated:
		c.onMessageCreated(ev.payload)
	case wsMessageUpdated:
		c.onMessageUpdated(ev.payload)
	case wsMessageDeleted:
		c.onMessageDeleted(ev.payload)
	case wsMessageButtonClicked:
		c.onMessageButtonClicked(ev.payload)
	case wsServerRoleCreated:
		c.onServerRoleCreated(ev.payload)
	case wsServerRoleDeleted:
		c.onServerRoleDeleted(ev.payload)
	case wsServerRoleUpdated:
		c.onServerRoleUpdated(ev.payload)
	case wsServerRoleOrderUpdate:
		c.onServerRoleOrderUpdated(ev.payload)
	case wsMessageReactionAdded:
		c.onMessageReactionAdded(ev.payload)
	case wsMessageReactionRemove:
		c.onMessageReactionRemoved(ev.payload)
	}
}

// onAuthenticated handles "user:authenticated". It decodes the payload one
// top-level field at a time rather than in a single json.Unmarshal call: the
// "user" field is by far the most important one to get right, since
// Client.User() being nil after this event fires is exactly the kind of bug
// that surfaces three call frames away as a nil pointer panic in unrelated
// code. If any other field (servers, channels, serverMembers, serverRoles)
// fails to decode — say, because the gateway added a field shape this SDK
// version doesn't know about — that failure is reported via OnError and
// skipped, but it no longer prevents Client.User() from being set and OnReady
// from firing.
func (c *Client) onAuthenticated(payload json.RawMessage) {
	var envelope struct {
		User          json.RawMessage `json:"user"`
		Servers       json.RawMessage `json:"servers"`
		Channels      json.RawMessage `json:"channels"`
		ServerMembers json.RawMessage `json:"serverMembers"`
		ServerRoles   json.RawMessage `json:"serverRoles"`
	}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding user:authenticated payload: %w", err))
		return
	}

	var rawSelf rawUser
	if err := json.Unmarshal(envelope.User, &rawSelf); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding user:authenticated \"user\" field: %w", err))
		return
	}
	c.mu.Lock()
	c.user = &ClientUser{User: newUser(c, rawSelf)}
	c.mu.Unlock()
	c.users.set(rawSelf.ID, c.user.User)

	var servers []rawServer
	if err := json.Unmarshal(envelope.Servers, &servers); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding user:authenticated \"servers\" field: %w", err))
	}
	for _, s := range servers {
		c.servers.setServer(s)
	}

	var channels []rawChannel
	if err := json.Unmarshal(envelope.Channels, &channels); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding user:authenticated \"channels\" field: %w", err))
	}
	for _, rawCh := range channels {
		ch := c.channels.setChannel(rawCh)
		if rawCh.ServerID != "" {
			if server, ok := c.servers.get(rawCh.ServerID); ok {
				server.channels.set(ch.ID, ch)
			}
		}
	}

	var members []rawServerMember
	if err := json.Unmarshal(envelope.ServerMembers, &members); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding user:authenticated \"serverMembers\" field: %w", err))
	}
	for _, m := range members {
		c.users.setUser(m.User)
		if server, ok := c.servers.get(m.ServerID); ok {
			member := newServerMember(c, m)
			server.members.set(member.ID, member)
		}
	}

	var roles []rawServerRole
	if err := json.Unmarshal(envelope.ServerRoles, &roles); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding user:authenticated \"serverRoles\" field: %w", err))
	}
	for _, r := range roles {
		if server, ok := c.servers.get(r.ServerID); ok {
			role := newServerRole(c, r)
			server.roles.set(role.ID, role)
		}
	}

	for _, fn := range c.snapshotHandlers().onReady {
		fn()
	}
}

func (c *Client) snapshotHandlers() eventHandlersSnapshot {
	c.handlers.mu.RLock()
	defer c.handlers.mu.RUnlock()
	return eventHandlersSnapshot{
		onReady:                  append([]func(){}, c.handlers.onReady...),
		onMessageCreate:          append([]func(*Message){}, c.handlers.onMessageCreate...),
		onMessageUpdate:          append([]func(*Message){}, c.handlers.onMessageUpdate...),
		onMessageDelete:          append([]func(MessageDeleteEvent){}, c.handlers.onMessageDelete...),
		onMessageReactionAdded:   append([]func(*Reaction, *User){}, c.handlers.onMessageReactionAdded...),
		onMessageReactionRemoved: append([]func(*Reaction, *User){}, c.handlers.onMessageReactionRemoved...),
		onMessageButtonClick:     append([]func(*MessageButton){}, c.handlers.onMessageButtonClick...),
		onServerMemberJoined:     append([]func(*ServerMember){}, c.handlers.onServerMemberJoined...),
		onServerMemberLeft:       append([]func(*ServerMember){}, c.handlers.onServerMemberLeft...),
		onServerMemberUpdated:    append([]func(*ServerMember){}, c.handlers.onServerMemberUpdated...),
		onServerJoined:           append([]func(*Server){}, c.handlers.onServerJoined...),
		onServerLeft:             append([]func(*Server){}, c.handlers.onServerLeft...),
		onServerChannelCreated:   append([]func(*ServerChannel){}, c.handlers.onServerChannelCreated...),
		onServerChannelUpdated:   append([]func(*ServerChannel){}, c.handlers.onServerChannelUpdated...),
		onServerChannelDeleted:   append([]func(ServerChannelDeleteEvent){}, c.handlers.onServerChannelDeleted...),
		onServerRoleCreated:      append([]func(*ServerRole){}, c.handlers.onServerRoleCreated...),
		onServerRoleUpdated:      append([]func(*ServerRole){}, c.handlers.onServerRoleUpdated...),
		onServerRoleDeleted:      append([]func(*ServerRole){}, c.handlers.onServerRoleDeleted...),
		onServerRoleOrderUpdated: append([]func(*Server, []*ServerRole){}, c.handlers.onServerRoleOrderUpdated...),
		onError:                  append([]func(error){}, c.handlers.onError...),
	}
}

// eventHandlersSnapshot is an immutable copy of the registered handler slices,
// taken under lock, so callbacks can be invoked without holding the lock (a
// handler registering another handler must not deadlock).
type eventHandlersSnapshot struct {
	onReady                  []func()
	onMessageCreate          []func(*Message)
	onMessageUpdate          []func(*Message)
	onMessageDelete          []func(MessageDeleteEvent)
	onMessageReactionAdded   []func(*Reaction, *User)
	onMessageReactionRemoved []func(*Reaction, *User)
	onMessageButtonClick     []func(*MessageButton)
	onServerMemberJoined     []func(*ServerMember)
	onServerMemberLeft       []func(*ServerMember)
	onServerMemberUpdated    []func(*ServerMember)
	onServerJoined           []func(*Server)
	onServerLeft             []func(*Server)
	onServerChannelCreated   []func(*ServerChannel)
	onServerChannelUpdated   []func(*ServerChannel)
	onServerChannelDeleted   []func(ServerChannelDeleteEvent)
	onServerRoleCreated      []func(*ServerRole)
	onServerRoleUpdated      []func(*ServerRole)
	onServerRoleDeleted      []func(*ServerRole)
	onServerRoleOrderUpdated []func(*Server, []*ServerRole)
	onError                  []func(error)
}

type serverMemberJoinedPayload struct {
	ServerID string          `json:"serverId"`
	Member   rawServerMember `json:"member"`
}

func (c *Client) onServerMemberJoined(payload json.RawMessage) {
	var p serverMemberJoinedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	server, ok := c.servers.get(p.ServerID)
	if !ok {
		return
	}
	c.users.setUser(p.Member.User)
	member := newServerMember(c, p.Member)
	server.members.set(member.ID, member)

	for _, fn := range c.snapshotHandlers().onServerMemberJoined {
		fn(member)
	}
}

type serverMemberUpdatedPayload struct {
	ServerID string `json:"serverId"`
	UserID   string `json:"userId"`
	Updated  struct {
		RoleIDs []string `json:"roleIds"`
	} `json:"updated"`
}

func (c *Client) onServerMemberUpdated(payload json.RawMessage) {
	var p serverMemberUpdatedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	server, ok := c.servers.get(p.ServerID)
	if !ok {
		return
	}
	member, ok := server.members.get(p.UserID)
	if !ok {
		return
	}
	if p.Updated.RoleIDs != nil {
		member.RoleIDs = p.Updated.RoleIDs
	}

	for _, fn := range c.snapshotHandlers().onServerMemberUpdated {
		fn(member)
	}
}

type serverJoinedPayload struct {
	Server   rawServer         `json:"server"`
	Members  []rawServerMember `json:"members"`
	Channels []rawChannel      `json:"channels"`
	Roles    []rawServerRole   `json:"roles"`
}

func (c *Client) onServerJoined(payload json.RawMessage) {
	var p serverJoinedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	server := c.servers.setServer(p.Server)

	for _, m := range p.Members {
		c.users.setUser(m.User)
		member := newServerMember(c, m)
		server.members.set(member.ID, member)
	}
	for _, r := range p.Roles {
		role := newServerRole(c, r)
		server.roles.set(role.ID, role)
	}
	for _, rawCh := range p.Channels {
		ch := c.channels.setChannel(rawCh)
		server.channels.set(ch.ID, ch)
	}

	for _, fn := range c.snapshotHandlers().onServerJoined {
		fn(server)
	}
}

type serverChannelCreatedPayload struct {
	ServerID string     `json:"serverId"`
	Channel  rawChannel `json:"channel"`
}

func (c *Client) onServerChannelCreated(payload json.RawMessage) {
	var p serverChannelCreatedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	ch := c.channels.setChannel(p.Channel)
	if server, ok := c.servers.get(p.ServerID); ok {
		server.channels.set(ch.ID, ch)
	}

	for _, fn := range c.snapshotHandlers().onServerChannelCreated {
		fn(ch)
	}
}

type serverChannelUpdatedPayload struct {
	ServerID  string     `json:"serverId"`
	ChannelID string     `json:"channelId"`
	Updated   rawChannel `json:"updated"`
}

func (c *Client) onServerChannelUpdated(payload json.RawMessage) {
	var p serverChannelUpdatedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	ch, ok := c.channels.get(p.ChannelID)
	if !ok {
		return
	}
	c.applyChannelUpdate(ch, payload)

	for _, fn := range c.snapshotHandlers().onServerChannelUpdated {
		fn(ch)
	}
}

// applyChannelUpdate merges only the fields present in the raw "updated" JSON
// object onto ch, so fields the server omitted are left untouched.
func (c *Client) applyChannelUpdate(ch *Channel, payload json.RawMessage) {
	var envelope struct {
		Updated map[string]json.RawMessage `json:"updated"`
	}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding server:channel_updated \"updated\" field: %w", err))
		return
	}
	if raw, ok := envelope.Updated["name"]; ok {
		_ = json.Unmarshal(raw, &ch.Name)
	}
	if raw, ok := envelope.Updated["permissions"]; ok {
		_ = json.Unmarshal(raw, &ch.Permissions)
	}
	if raw, ok := envelope.Updated["categoryId"]; ok {
		_ = json.Unmarshal(raw, &ch.CategoryID)
	}
	if raw, ok := envelope.Updated["order"]; ok {
		_ = json.Unmarshal(raw, &ch.Order)
	}
}

type serverChannelDeletedPayload struct {
	ServerID  string `json:"serverId"`
	ChannelID string `json:"channelId"`
}

func (c *Client) onServerChannelDeleted(payload json.RawMessage) {
	var p serverChannelDeletedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	if !c.channels.has(p.ChannelID) {
		return
	}
	c.channels.delete(p.ChannelID)
	if server, ok := c.servers.get(p.ServerID); ok {
		server.channels.delete(p.ChannelID)
	}

	ev := ServerChannelDeleteEvent{ServerID: p.ServerID, ChannelID: p.ChannelID}
	for _, fn := range c.snapshotHandlers().onServerChannelDeleted {
		fn(ev)
	}
}

type serverLeftPayload struct {
	ServerID string `json:"serverId"`
}

func (c *Client) onServerLeft(payload json.RawMessage) {
	var p serverLeftPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	server, ok := c.servers.get(p.ServerID)
	if !ok {
		return
	}

	for _, fn := range c.snapshotHandlers().onServerLeft {
		fn(server)
	}

	c.servers.delete(p.ServerID)
	for _, ch := range c.channels.values() {
		if ch.ServerID == p.ServerID {
			c.channels.delete(ch.ID)
		}
	}
	server.members.clear()
}

type serverMemberLeftPayload struct {
	UserID   string `json:"userId"`
	ServerID string `json:"serverId"`
}

func (c *Client) onServerMemberLeft(payload json.RawMessage) {
	var p serverMemberLeftPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	server, ok := c.servers.get(p.ServerID)
	if !ok {
		return
	}
	member, ok := server.members.get(p.UserID)
	if !ok {
		return
	}

	for _, fn := range c.snapshotHandlers().onServerMemberLeft {
		fn(member)
	}
	server.members.delete(p.UserID)
}

type messageCreatedPayload struct {
	Message rawMessage `json:"message"`
}

func (c *Client) onMessageCreated(payload json.RawMessage) {
	var p messageCreatedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	msg := c.messages.setMessage(p.Message)

	for _, fn := range c.snapshotHandlers().onMessageCreate {
		fn(msg)
	}
}

type messageUpdatedPayload struct {
	ChannelID string     `json:"channelId"`
	MessageID string     `json:"messageId"`
	Updated   rawMessage `json:"updated"`
}

func (c *Client) onMessageUpdated(payload json.RawMessage) {
	var p messageUpdatedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	msg, ok := c.messages.get(p.MessageID)
	if !ok {
		return
	}
	msg.update(p.Updated)
	c.messages.set(msg.ID, msg)

	for _, fn := range c.snapshotHandlers().onMessageUpdate {
		fn(msg)
	}
}

type messageDeletedPayload struct {
	ChannelID string `json:"channelId"`
	MessageID string `json:"messageId"`
}

func (c *Client) onMessageDeleted(payload json.RawMessage) {
	var p messageDeletedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	c.messages.delete(p.MessageID)

	ev := MessageDeleteEvent{MessageID: p.MessageID, ChannelID: p.ChannelID}
	for _, fn := range c.snapshotHandlers().onMessageDelete {
		fn(ev)
	}
}

func (c *Client) onMessageButtonClicked(payload json.RawMessage) {
	var p messageButtonClickPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	button := newMessageButton(c, p)

	for _, fn := range c.snapshotHandlers().onMessageButtonClick {
		fn(button)
	}
}

func (c *Client) onServerRoleCreated(payload json.RawMessage) {
	var p rawServerRole
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	server, ok := c.servers.get(p.ServerID)
	if !ok {
		return
	}
	role := newServerRole(c, p)
	server.roles.set(role.ID, role)

	for _, fn := range c.snapshotHandlers().onServerRoleCreated {
		fn(role)
	}
}

type serverRoleDeletedPayload struct {
	ServerID string `json:"serverId"`
	RoleID   string `json:"roleId"`
}

func (c *Client) onServerRoleDeleted(payload json.RawMessage) {
	var p serverRoleDeletedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	server, ok := c.servers.get(p.ServerID)
	if !ok {
		return
	}
	role, ok := server.roles.get(p.RoleID)
	if !ok {
		return
	}
	server.roles.delete(p.RoleID)

	for _, fn := range c.snapshotHandlers().onServerRoleDeleted {
		fn(role)
	}
}

type serverRoleUpdatedPayload struct {
	ServerID string `json:"serverId"`
	RoleID   string `json:"roleId"`
}

func (c *Client) onServerRoleUpdated(payload json.RawMessage) {
	var p serverRoleUpdatedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	server, ok := c.servers.get(p.ServerID)
	if !ok {
		return
	}
	role, ok := server.roles.get(p.RoleID)
	if !ok {
		return
	}
	c.applyRoleUpdate(role, payload)

	for _, fn := range c.snapshotHandlers().onServerRoleUpdated {
		fn(role)
	}
}

func (c *Client) applyRoleUpdate(role *ServerRole, payload json.RawMessage) {
	var envelope struct {
		Updated map[string]json.RawMessage `json:"updated"`
	}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding server:role_updated \"updated\" field: %w", err))
		return
	}
	if raw, ok := envelope.Updated["name"]; ok {
		_ = json.Unmarshal(raw, &role.Name)
	}
	if raw, ok := envelope.Updated["permissions"]; ok {
		_ = json.Unmarshal(raw, &role.Permissions)
	}
	if raw, ok := envelope.Updated["hexColor"]; ok {
		_ = json.Unmarshal(raw, &role.HexColor)
	}
	if raw, ok := envelope.Updated["order"]; ok {
		_ = json.Unmarshal(raw, &role.Order)
	}
}

type serverRoleOrderUpdatedPayload struct {
	ServerID string   `json:"serverId"`
	RoleIDs  []string `json:"roleIds"`
}

func (c *Client) onServerRoleOrderUpdated(payload json.RawMessage) {
	var p serverRoleOrderUpdatedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	server, ok := c.servers.get(p.ServerID)
	if !ok {
		return
	}
	for i, roleID := range p.RoleIDs {
		if role, ok := server.roles.get(roleID); ok {
			role.Order = i + 1
		}
	}

	for _, fn := range c.snapshotHandlers().onServerRoleOrderUpdated {
		fn(server, server.roles.values())
	}
}

func (c *Client) onMessageReactionAdded(payload json.RawMessage) {
	var p reactionAddedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	reaction := newReactionFromAdded(c, p)
	user, _ := c.users.get(p.ReactedByUserID)

	for _, fn := range c.snapshotHandlers().onMessageReactionAdded {
		fn(reaction, user)
	}
}

func (c *Client) onMessageReactionRemoved(payload json.RawMessage) {
	var p reactionRemovedPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		c.reportError(fmt.Errorf("nerimity: decoding gateway event payload: %w", err))
		return
	}
	reaction := newReactionFromRemoved(c, p)
	user, _ := c.users.get(p.ReactionRemovedByUser)

	for _, fn := range c.snapshotHandlers().onMessageReactionRemoved {
		fn(reaction, user)
	}
}
