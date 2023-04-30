package commands

import (
	"context"
	"fmt"
	"log"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"github.com/SushiWaUmai/prince/db"
)

func init() {
	createCommand("clearrepeat", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) {
		// Delete the message
		affected := db.ClearRepeatedMessage(messageEvent.Info.Chat.String())

		// Reply
		_, err := client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
			Conversation: proto.String(fmt.Sprintf("Deleted %d", affected)),
		})

		if err != nil {
			log.Println(err)
		}
	})
}
