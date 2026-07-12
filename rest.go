package nerimity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIError is returned when the Nerimity REST API responds with a non-2xx
// status. Body holds the raw response body for inspection.
type APIError struct {
	Status int
	Body   string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("nerimity: API returned %d: %s", e.Status, e.Body)
}

// doJSON performs a JSON request. If body is non-nil it is marshalled as the
// request body; if out is non-nil the response body is unmarshalled into it.
// authToken, when non-empty, is sent as the Authorization header (Nerimity uses
// the raw token, with no "Bearer" prefix).
func (c *Client) doJSON(ctx context.Context, method, url string, body, out any, authToken string) error {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("nerimity: encoding request body: %w", err)
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return fmt.Errorf("nerimity: building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("nerimity: could not connect to server: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{Status: resp.StatusCode, Body: string(raw)}
	}
	if out != nil && len(raw) > 0 {
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("nerimity: decoding response: %w", err)
		}
	}
	return nil
}

func (c *Client) messagesURL(channelID string) string {
	return c.apiBase + "/channels/" + channelID + "/messages"
}

func (c *Client) messageURL(channelID, messageID string) string {
	return c.messagesURL(channelID) + "/" + messageID
}

type postMessageBody struct {
	Content           string         `json:"content"`
	NerimityCdnFileID string         `json:"nerimityCdnFileId,omitempty"`
	HTMLEmbed         string         `json:"htmlEmbed,omitempty"`
	Buttons           []ButtonOption `json:"buttons,omitempty"`
	Silent            bool           `json:"silent,omitempty"`
	MentionReplies    bool           `json:"mentionReplies,omitempty"`
	ReplyToMessageIDs []string       `json:"replyToMessageIds,omitempty"`
}

func (c *Client) postMessage(ctx context.Context, channelID, content string, opts MessageOptions) (*Message, error) {
	body := postMessageBody{
		Content:           content,
		NerimityCdnFileID: opts.NerimityCdnFileID,
		HTMLEmbed:         opts.HTMLEmbed,
		Buttons:           opts.Buttons,
		Silent:            opts.Silent,
		MentionReplies:    opts.MentionReplies,
		ReplyToMessageIDs: opts.ReplyToMessageIDs,
	}
	var raw rawMessage
	if err := c.doJSON(ctx, http.MethodPost, c.messagesURL(channelID), body, &raw, c.token); err != nil {
		return nil, err
	}
	return newMessage(c, raw), nil
}

type editMessageBody struct {
	Content   string         `json:"content"`
	HTMLEmbed string         `json:"htmlEmbed,omitempty"`
	Buttons   []ButtonOption `json:"buttons,omitempty"`
}

func (c *Client) editMessage(ctx context.Context, channelID, messageID, content string, opts EditOptions) (*Message, error) {
	body := editMessageBody{Content: content, HTMLEmbed: opts.HTMLEmbed, Buttons: opts.Buttons}
	var raw rawMessage
	if err := c.doJSON(ctx, http.MethodPatch, c.messageURL(channelID, messageID), body, &raw, c.token); err != nil {
		return nil, err
	}
	return newMessage(c, raw), nil
}

func (c *Client) deleteMessage(ctx context.Context, channelID, messageID string) error {
	return c.doJSON(ctx, http.MethodDelete, c.messageURL(channelID, messageID), nil, nil, c.token)
}

func (c *Client) fetchMessage(ctx context.Context, channelID, messageID string) (*Message, error) {
	var raw rawMessage
	if err := c.doJSON(ctx, http.MethodGet, c.messageURL(channelID, messageID), nil, &raw, c.token); err != nil {
		return nil, err
	}
	return newMessage(c, raw), nil
}

func (c *Client) buttonCallback(ctx context.Context, channelID, messageID, buttonID, userID string, resp ButtonResponse) error {
	url := c.messageURL(channelID, messageID) + "/buttons/" + buttonID + "/callback"
	body := map[string]any{"userId": userID}
	if resp.Content != "" {
		body["content"] = resp.Content
	}
	if len(resp.Components) > 0 {
		body["components"] = resp.Components
	}
	if resp.Title != "" {
		body["title"] = resp.Title
	}
	if resp.ButtonLabel != "" {
		body["buttonLabel"] = resp.ButtonLabel
	}
	return c.doJSON(ctx, http.MethodPost, url, body, nil, c.token)
}

func (c *Client) banMember(ctx context.Context, serverID, userID, reason string) error {
	url := c.apiBase + "/servers/" + serverID + "/bans/" + userID
	return c.doJSON(ctx, http.MethodPost, url, map[string]any{"reason": reason}, nil, c.token)
}

func (c *Client) unbanMember(ctx context.Context, serverID, userID string) error {
	url := c.apiBase + "/servers/" + serverID + "/bans/" + userID
	return c.doJSON(ctx, http.MethodDelete, url, nil, nil, c.token)
}

func (c *Client) kickMember(ctx context.Context, serverID, userID string) error {
	url := c.apiBase + "/servers/" + serverID + "/members/" + userID + "/kick"
	return c.doJSON(ctx, http.MethodDelete, url, nil, nil, c.token)
}
