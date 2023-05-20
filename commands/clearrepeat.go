package commands

import (
	"fmt"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"github.com/SushiWaUmai/prince/db"
)

func init() {
	createCommand("clearrepeat", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		// Delete the message
		affected := db.ClearRepeatedMessage(messageEvent.Info.Chat.String())

		// Reply
		response := &waProto.Message{
			Conversation: proto.String(fmt.Sprintf("Deleted %d", affected)),
		}

		return response, nil
	})
}
