package commands

import (
	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	utils.CreateCommand("ping", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
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
