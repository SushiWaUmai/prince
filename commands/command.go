package commands

import (
	"context"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var CommandList []Command

type Command struct {
	Name    string
	Execute func(client *whatsmeow.Client, messageEvent *events.Message, args []string)
}

func createCommand(name string, execute func(client *whatsmeow.Client, messageEvent *events.Message, args []string)) {
	CommandList = append(CommandList, Command{
		Name:    name,
		Execute: execute,
	})
}

func init() {
	createCommand("ping", func(client *whatsmeow.Client, messageEvent *events.Message, args []string) {
		reply := "pong!"
		if len(args) > 0 {
			reply = args[0]
		}

		client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
			Conversation: proto.String(reply),
		})
	})
}
