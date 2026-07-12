# Webhooks

`WebhookBuilder` sends messages to a channel without a bot account or
connection; just an HTTP POST per message.

## From a webhook URL

```go
webhook, err := nerimity.NewWebhook("https://nerimity.com/api/webhooks/1234567890/the-webhook-token")
if err != nil {
	log.Fatal(err)
}
```

## From channel ID and token

```go
webhook, err := nerimity.NewWebhookFromParts("1234567890", "the-webhook-token")
```

## Sending

```go
webhook.SetUsername("Deploy Bot").
	SetAvatar("https://example.com/avatar.png")

if err := webhook.Send(context.Background(), "Deployment finished."); err != nil {
	log.Fatal(err)
}
```

`SetUsername` and `SetAvatar` are optional and override the webhook's default
display name and avatar for that message; both return the builder so calls
chain. Reuse the same `*WebhookBuilder` for multiple sends; it's not
single-use.

Webhooks are independent of `Client`: you don't need to call `Login` or hold a
`Client` at all to send one. See [`examples/webhook-sender`](../examples/webhook-sender)
for a complete program.
