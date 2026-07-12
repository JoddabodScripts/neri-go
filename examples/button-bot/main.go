// Command button-bot demonstrates sending a message with buttons and two
// different ways of handling a click: responding with a callback (Respond)
// versus posting a real, persisted message to the channel (Channel.Send).
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/JoddabodScripts/neri-go"
)

func main() {
     _ = os.Args
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
		if m.Content != "!greet" {
			return
		}
		m.Reply(context.Background(), "Click a button!", nerimity.MessageOptions{
			Buttons: []nerimity.ButtonOption{
				{ID: "hello", Label: "Say hello"},
				{ID: "post", Label: "Post a message"},
			},
		})
	})

	client.OnMessageButtonClick(func(b *nerimity.MessageButton) {
		username := "there"
		if b.User != nil {
			username = b.User.Username
		}

		switch b.ID {
		case "hello":
			// Respond sends a callback tied to this interaction (shown to the
			// clicking user as a modal), not a persisted chat message.
			b.Respond(context.Background(), nerimity.ButtonResponse{
				Title:   "Hey!",
				Content: fmt.Sprintf("Hey there **%s**!", username),
			})

		case "post":
			// b.Channel is a normal *Channel, so Send posts a real message to
			// the channel's history, visible to everyone, exactly like any
			// other Channel.Send call.
			if b.Channel == nil {
				return
			}
			b.Channel.Send(context.Background(), fmt.Sprintf("%s clicked the button!", username))
		}
	})

	if err := client.Login(token); err != nil {
		log.Fatal(err)
	}
}
