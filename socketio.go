package nerimity

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// This file implements the minimal subset of the Engine.IO v4 / Socket.IO v5
// protocol that the Nerimity gateway uses. The reference client
// (socket.io-client, configured with transports: ["websocket"]) skips HTTP
// long-polling and connects straight over WebSocket, which is exactly what we
// do here.
//
// Frame format (text WebSocket frames):
//   Engine.IO packet type is the first character:
//     0 open   1 close   2 ping   3 pong   4 message   6 noop
//   For a message ("4..."), the next character is the Socket.IO packet type:
//     0 CONNECT   1 DISCONNECT   2 EVENT   3 ACK   4 CONNECT_ERROR
//   An EVENT ("42...") is followed by an optional namespace, an optional ack
//   id, then a JSON array: ["eventName", arg].

const (
	engineOpen    = '0'
	engineClose   = '1'
	enginePing    = '2'
	enginePong    = '3'
	engineMessage = '4'

	socketConnect      = '0'
	socketDisconnect   = '1'
	socketEvent        = '2'
	socketConnectError = '4'
)

// socket is a single Socket.IO connection to the gateway. It is not reused
// across reconnects; the Client creates a fresh one each time.
type socket struct {
	conn    *websocket.Conn
	writeMu sync.Mutex

	pingInterval time.Duration
	pingTimeout  time.Duration
}

type engineHandshake struct {
	SID          string `json:"sid"`
	PingInterval int    `json:"pingInterval"`
	PingTimeout  int    `json:"pingTimeout"`
}

// gatewayURL converts an http(s) base URL into the Socket.IO WebSocket
// endpoint.
func gatewayURL(base string) string {
	u := base
	switch {
	case strings.HasPrefix(u, "https://"):
		u = "wss://" + strings.TrimPrefix(u, "https://")
	case strings.HasPrefix(u, "http://"):
		u = "ws://" + strings.TrimPrefix(u, "http://")
	}
	u = strings.TrimSuffix(u, "/")
	return u + "/socket.io/?EIO=4&transport=websocket"
}

// connect dials the gateway, completes the Engine.IO and Socket.IO handshakes,
// and returns once the connection is ready to send and receive events.
func (s *socket) connect(ctx context.Context, base string) error {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, gatewayURL(base), nil)
	if err != nil {
		return fmt.Errorf("nerimity: dialing gateway: %w", err)
	}
	s.conn = conn

	// First frame must be the Engine.IO open packet.
	frame, err := s.readFrame()
	if err != nil {
		conn.Close()
		return fmt.Errorf("nerimity: reading handshake: %w", err)
	}
	if len(frame) == 0 || frame[0] != engineOpen {
		conn.Close()
		return fmt.Errorf("nerimity: unexpected handshake frame: %q", frame)
	}
	var hs engineHandshake
	if err := json.Unmarshal([]byte(frame[1:]), &hs); err != nil {
		conn.Close()
		return fmt.Errorf("nerimity: decoding handshake: %w", err)
	}
	s.pingInterval = durOrDefault(hs.PingInterval, 25000)
	s.pingTimeout = durOrDefault(hs.PingTimeout, 20000)

	// Connect to the default Socket.IO namespace.
	if err := s.write(string(engineMessage) + string(socketConnect)); err != nil {
		conn.Close()
		return err
	}
	// Expect the CONNECT acknowledgement ("40{...}").
	for {
		frame, err := s.readFrame()
		if err != nil {
			conn.Close()
			return fmt.Errorf("nerimity: reading connect ack: %w", err)
		}
		if len(frame) >= 2 && frame[0] == engineMessage && frame[1] == socketConnect {
			return nil
		}
		if len(frame) >= 2 && frame[0] == engineMessage && frame[1] == socketConnectError {
			conn.Close()
			return fmt.Errorf("nerimity: gateway rejected connection: %s", frame[2:])
		}
		if len(frame) >= 1 && frame[0] == enginePing {
			_ = s.write(string(enginePong))
		}
	}
}

func durOrDefault(ms, def int) time.Duration {
	if ms <= 0 {
		ms = def
	}
	return time.Duration(ms) * time.Millisecond
}

// read blocks until the next Socket.IO EVENT arrives, transparently answering
// Engine.IO pings and ignoring other control frames. It returns the event name
// and its first argument (the payload object), or an error when the connection
// closes.
func (s *socket) read() (name string, payload json.RawMessage, err error) {
	for {
		frame, err := s.readFrame()
		if err != nil {
			return "", nil, err
		}
		if len(frame) == 0 {
			continue
		}
		switch frame[0] {
		case enginePing:
			if err := s.write(string(enginePong)); err != nil {
				return "", nil, err
			}
		case engineClose:
			return "", nil, fmt.Errorf("nerimity: gateway closed the connection")
		case engineMessage:
			if len(frame) < 2 {
				continue
			}
			switch frame[1] {
			case socketEvent:
				name, payload, ok := parseEvent(frame[2:])
				if ok {
					return name, payload, nil
				}
			case socketDisconnect:
				return "", nil, fmt.Errorf("nerimity: gateway disconnected the session")
			case socketConnectError:
				return "", nil, fmt.Errorf("nerimity: gateway connect error: %s", frame[2:])
			}
		}
	}
}

// parseEvent decodes the body of a Socket.IO EVENT packet (everything after the
// "42" prefix) into an event name and its first argument.
func parseEvent(body string) (name string, payload json.RawMessage, ok bool) {
	// Skip an optional namespace ("/foo,") and an optional numeric ack id; the
	// JSON array always starts at the first '['.
	idx := strings.IndexByte(body, '[')
	if idx < 0 {
		return "", nil, false
	}
	var args []json.RawMessage
	if err := json.Unmarshal([]byte(body[idx:]), &args); err != nil || len(args) == 0 {
		return "", nil, false
	}
	if err := json.Unmarshal(args[0], &name); err != nil {
		return "", nil, false
	}
	if len(args) > 1 {
		payload = args[1]
	}
	return name, payload, true
}

// emit sends a Socket.IO EVENT to the gateway.
func (s *socket) emit(event string, arg any) error {
	var arr []any
	if arg == nil {
		arr = []any{event}
	} else {
		arr = []any{event, arg}
	}
	data, err := json.Marshal(arr)
	if err != nil {
		return err
	}
	return s.write(string(engineMessage) + string(socketEvent) + string(data))
}

// readFrame reads one text frame, refreshing the read deadline based on the
// negotiated ping timing so a dead connection is detected.
func (s *socket) readFrame() (string, error) {
	if s.pingInterval > 0 {
		_ = s.conn.SetReadDeadline(time.Now().Add(s.pingInterval + s.pingTimeout))
	}
	_, data, err := s.conn.ReadMessage()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *socket) write(data string) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	return s.conn.WriteMessage(websocket.TextMessage, []byte(data))
}

func (s *socket) close() error {
	if s.conn == nil {
		return nil
	}
	return s.conn.Close()
}
