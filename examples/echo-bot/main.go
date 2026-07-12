// Command echo-bot is a minimal Nerimity bot: it replies "Pong!" to "!ping" and
// echoes back anything else a user sends.
package main

import (
	"context"
	"log"
	"os"

	"github.com/JoddabodScripts/neri-go"
)

func main() {
	token := os.Getenv("NERIMITY_TOKEN")
	if token == "" {
		log.Fatal("set NERIMITY_TOKEN to your bot's token")
	}

	client := nerimity.New(nerimity.Options{})

	client.OnReady(func() {
		log.Printf("connected as %s", client.User().Username)
	})

	client.OnMessageCreate(func(m *nerimity.Message) {
		if m.User == nil || m.User.ID == client.User().ID {
			return
		}
		ctx := context.Background()
		if m.Content == "!ping" {
			m.Reply(ctx, "Pong!")
			return
		}
		if m.Content != "" {
			m.Channel.Send(ctx, m.Content)
		}
	})

	if err := client.Login(token); err != nil {
		log.Fatal(err)
	}
}
