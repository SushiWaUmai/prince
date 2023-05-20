package commands

import (
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	createCommand("ping", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		reply := "pong!"
		if len(args) > 0 {
			reply = args[0]
		}

		response := &waProto.Message{
			Conversation: &reply,
		}
		return response, nil
	})
}
