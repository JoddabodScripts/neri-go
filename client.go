package nerimity

import (
	"context"
	"encoding/json"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Default endpoints and cache sizes.
const (
	defaultAPIURL            = "https://nerimity.com"
	defaultWSURL             = "https://nerimity.com"
	defaultMessageCacheLimit = 1000
)

// Socket.IO event names sent by the client to the gateway.
const (
	eventAuthenticate   = "user:authenticate"
	eventUpdateActivity = "user:update_activity"
)

// Options configures a Client. The zero value is valid and uses Nerimity's
// production endpoints with a 1000-message cache.
type Options struct {
	// WSURLOverride overrides the gateway (WebSocket) base URL. Defaults to
	// https://nerimity.com.
	WSURLOverride string
	// APIURLOverride overrides the REST API base URL. Defaults to
	// https://nerimity.com.
	APIURLOverride string
	// MessageCacheLimit caps the number of cached messages (LRU). Defaults to
	// 1000. Set to -1 for unbounded.
	MessageCacheLimit int
	// HTTPClient is used for all REST and CDN requests. Defaults to a client
	// with a 30-second timeout.
	HTTPClient *http.Client
}

// Client is a connection to Nerimity. Register event handlers with the On*
// methods, then call Login to connect. A Client must not be copied.
type Client struct {
	apiBase    string
	wsURL      string
	httpClient *http.Client

	token     string
	reconnect bool

	users    *userStore
	servers  *serverStore
	channels *channelStore
	messages *messageStore

	handlers *eventHandlers

	mu     sync.RWMutex
	user   *ClientUser
	sock   *socket
	cancel context.CancelFunc

	inbound chan incomingEvent
}

type incomingEvent struct {
	name    string
	payload json.RawMessage
}

// New creates a Client with the given options.
func New(opts Options) *Client {
	apiURL := opts.APIURLOverride
	if apiURL == "" {
		apiURL = defaultAPIURL
	}
	wsURL := opts.WSURLOverride
	if wsURL == "" {
		wsURL = defaultWSURL
	}
	limit := opts.MessageCacheLimit
	switch {
	case limit == 0:
		limit = defaultMessageCacheLimit
	case limit < 0:
		limit = 0 // unbounded
	}
	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	c := &Client{
		apiBase:    apiURL + "/api",
		wsURL:      wsURL,
		httpClient: httpClient,
		reconnect:  true,
		handlers:   &eventHandlers{},
		inbound:    make(chan incomingEvent, 128),
	}
	c.users = &userStore{cache: newCache[*User](0), client: c}
	c.servers = &serverStore{cache: newCache[*Server](0), client: c}
	c.channels = &channelStore{cache: newCache[*Channel](0), client: c}
	c.messages = &messageStore{cache: newCache[*Message](limit), client: c}
	return c
}

// User returns the bot's own user. It is nil until the Ready event fires.
func (c *Client) User() *ClientUser {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.user
}

func (c *Client) selfID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.user != nil {
		return c.user.ID
	}
	return ""
}

// Servers returns a snapshot of the cached servers.
func (c *Client) Servers() []*Server { return c.servers.values() }

// Server returns the cached server with the given ID, or nil.
func (c *Client) Server(id string) *Server { s, _ := c.servers.get(id); return s }

// Channel returns the cached channel with the given ID, or nil.
func (c *Client) Channel(id string) *Channel { ch, _ := c.channels.get(id); return ch }

// GetUser returns the cached user with the given ID, or nil.
func (c *Client) GetUser(id string) *User { u, _ := c.users.get(id); return u }

// SetActivity updates the bot's activity. Pass nil to clear it. Requires an
// active connection.
func (c *Client) SetActivity(activity *Activity) error {
	c.mu.RLock()
	sock := c.sock
	c.mu.RUnlock()
	if sock == nil {
		return nil
	}
	if activity == nil {
		return sock.emit(eventUpdateActivity, nil)
	}
	return sock.emit(eventUpdateActivity, activity)
}

// Activity is the bot's rich-presence activity, set with SetActivity.
type Activity struct {
	Action    string `json:"action"`
	Name      string `json:"name"`
	StartedAt int64  `json:"startedAt"`
	EndsAt    int64  `json:"endsAt,omitempty"`
	ImgSrc    string `json:"imgSrc,omitempty"`
	Title     string `json:"title,omitempty"`
	Subtitle  string `json:"subtitle,omitempty"`
	Link      string `json:"link,omitempty"`
}

// Login connects to Nerimity with the given bot token and blocks, dispatching
// events to registered handlers and reconnecting automatically, until Close is
// called or the context is cancelled. It is equivalent to LoginWithContext with
// a background context.
func (c *Client) Login(token string) error {
	return c.LoginWithContext(context.Background(), token)
}

// LoginWithContext is Login with a caller-supplied context. When ctx is
// cancelled the connection is torn down and the returned error is ctx.Err().
func (c *Client) LoginWithContext(ctx context.Context, token string) error {
	c.token = token
	ctx, cancel := context.WithCancel(ctx)
	c.mu.Lock()
	c.cancel = cancel
	c.mu.Unlock()
	defer cancel()

	go c.dispatchLoop(ctx)
	return c.connectLoop(ctx)
}

// Close disconnects the client and causes Login to return.
func (c *Client) Close() {
	c.mu.Lock()
	if c.cancel != nil {
		c.cancel()
	}
	sock := c.sock
	c.mu.Unlock()
	if sock != nil {
		sock.close()
	}
}

func (c *Client) connectLoop(ctx context.Context) error {
	attempt := 0
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		sock := &socket{}
		err := sock.connect(ctx, c.wsURL)
		if err == nil {
			attempt = 0
			c.mu.Lock()
			c.sock = sock
			c.mu.Unlock()
			err = c.runSession(ctx, sock)
		}
		sock.close()

		if err := ctx.Err(); err != nil {
			return err
		}
		if !c.reconnect {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoffDelay(attempt)):
		}
		attempt++
	}
}

// runSession authenticates and pumps gateway events into the inbound channel
// until the connection fails.
func (c *Client) runSession(ctx context.Context, sock *socket) error {
	if err := sock.emit(eventAuthenticate, map[string]string{"token": c.token}); err != nil {
		return err
	}
	for {
		name, payload, err := sock.read()
		if err != nil {
			return err
		}
		select {
		case c.inbound <- incomingEvent{name: name, payload: payload}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *Client) dispatchLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-c.inbound:
			c.handleEvent(ev)
		}
	}
}

// backoffDelay mirrors socket.io-client's reconnection timing: base 1s,
// doubling, capped at 5s, with +/-50% jitter.
func backoffDelay(attempt int) time.Duration {
	const base = float64(time.Second)
	const max = float64(5 * time.Second)
	d := base * math.Pow(2, float64(attempt))
	if d > max {
		d = max
	}
	jitter := d * 0.5 * (rand.Float64()*2 - 1)
	return time.Duration(d + jitter)
}

// ---- cache stores ----

type userStore struct {
	*cache[*User]
	client *Client
}

func (s *userStore) setUser(raw rawUser) *User {
	u := newUser(s.client, raw)
	s.set(raw.ID, u)
	return u
}

type serverStore struct {
	*cache[*Server]
	client *Client
}

func (s *serverStore) setServer(raw rawServer) *Server {
	srv := newServer(s.client, raw)
	s.set(raw.ID, srv)
	return srv
}

type channelStore struct {
	*cache[*Channel]
	client *Client
}

func (s *channelStore) setChannel(raw rawChannel) *Channel {
	ch := newChannel(s.client, raw)
	s.set(raw.ID, ch)
	return ch
}

type messageStore struct {
	*cache[*Message]
	client *Client
}

func (s *messageStore) setMessage(raw rawMessage) *Message {
	m := newMessage(s.client, raw)
	s.set(raw.ID, m)
	return m
}
