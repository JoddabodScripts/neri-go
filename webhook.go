package nerimity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

// webhookBaseURL is where webhook messages are POSTed. Webhooks do not use a
// bot token or the API override.
const webhookBaseURL = "https://nerimity.com/api/webhooks"

var webhookURLRegex = regexp.MustCompile(`/webhooks/(\d+)/(.+)$`)

// WebhookBuilder sends messages to a Nerimity channel through a webhook,
// without a bot account. Construct it from a webhook URL with NewWebhook, or
// from a channel ID and token with NewWebhookFromParts.
type WebhookBuilder struct {
	channelID  string
	token      string
	username   string
	avatarURL  string
	httpClient *http.Client
}

// NewWebhook creates a WebhookBuilder from a full webhook URL of the form
// "https://nerimity.com/api/webhooks/{channelId}/{token}".
func NewWebhook(webhookURL string) (*WebhookBuilder, error) {
	m := webhookURLRegex.FindStringSubmatch(webhookURL)
	if m == nil {
		return nil, fmt.Errorf("nerimity: invalid webhook URL: %q", webhookURL)
	}
	return &WebhookBuilder{channelID: m[1], token: m[2], httpClient: http.DefaultClient}, nil
}

// NewWebhookFromParts creates a WebhookBuilder from a channel ID and webhook
// token.
func NewWebhookFromParts(channelID, token string) (*WebhookBuilder, error) {
	if channelID == "" || token == "" {
		return nil, fmt.Errorf("nerimity: webhook requires both channelId and token")
	}
	return &WebhookBuilder{channelID: channelID, token: token, httpClient: http.DefaultClient}, nil
}

// SetUsername overrides the display name for webhook messages. Returns the
// builder for chaining.
func (w *WebhookBuilder) SetUsername(username string) *WebhookBuilder {
	w.username = username
	return w
}

// SetAvatar overrides the avatar (by image URL) for webhook messages. Returns
// the builder for chaining.
func (w *WebhookBuilder) SetAvatar(avatarURL string) *WebhookBuilder {
	w.avatarURL = avatarURL
	return w
}

// SetHTTPClient sets the HTTP client used to send. Optional; defaults to
// http.DefaultClient.
func (w *WebhookBuilder) SetHTTPClient(c *http.Client) *WebhookBuilder {
	w.httpClient = c
	return w
}

type webhookBody struct {
	Content   string `json:"content"`
	Username  string `json:"username,omitempty"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

// Send posts a message through the webhook.
func (w *WebhookBuilder) Send(ctx context.Context, content string) error {
	body, err := json.Marshal(webhookBody{Content: content, Username: w.username, AvatarURL: w.avatarURL})
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s/%s", webhookBaseURL, w.channelID, w.token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("nerimity: sending webhook: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		return &APIError{Status: resp.StatusCode, Body: string(raw)}
	}
	return nil
}
