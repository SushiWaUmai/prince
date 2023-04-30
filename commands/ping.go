package commands

import (
	"context"
	"log"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	createCommand("ping", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) {
		reply := "pong!"
		if len(args) > 0 {
			reply = args[0]
		}

		_, err := client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
			Conversation: &reply,
		})

		if err != nil {
			log.Println(err)
		}
	})
}
