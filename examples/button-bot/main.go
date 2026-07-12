// Command button-bot demonstrates sending a message with a button and
// responding when a user clicks it.
package main

import (
	"context"
	"fmt"
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
		if m.Content != "!greet" {
			return
		}
		m.Reply(context.Background(), "Click the button!", nerimity.MessageOptions{
			Buttons: []nerimity.ButtonOption{
				{ID: "hello", Label: "Say hello"},
			},
		})
	})

	client.OnMessageButtonClick(func(b *nerimity.MessageButton) {
		if b.ID != "hello" {
			return
		}
		username := "there"
		if b.User != nil {
			username = b.User.Username
		}
		b.Respond(context.Background(), nerimity.ButtonResponse{
			Title:   "Hey!",
			Content: fmt.Sprintf("Hey there **%s**!", username),
		})
	})

	if err := client.Login(token); err != nil {
		log.Fatal(err)
	}
}
