// Command webhook-sender posts a single message through a Nerimity webhook,
// without a bot account.
package main

import (
	"context"
	"log"
	"os"

	"github.com/JoddabodScripts/neri-go"
)

func main() {
	webhookURL := os.Getenv("NERIMITY_WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("set NERIMITY_WEBHOOK_URL to your webhook's URL")
	}

	webhook, err := nerimity.NewWebhook(webhookURL)
	if err != nil {
		log.Fatal(err)
	}

	webhook.SetUsername("Deploy Bot").
		SetAvatar("https://example.com/avatar.png")

	if err := webhook.Send(context.Background(), "Deployment finished successfully."); err != nil {
		log.Fatal(err)
	}
}
