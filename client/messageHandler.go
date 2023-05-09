package client

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/SushiWaUmai/prince/commands"
	"github.com/SushiWaUmai/prince/db"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func (client *PrinceClient) handleMessage(e *events.Message) {
	if !e.Info.IsFromMe {
		return
	}

	content, ctx := commands.GetTextContext(e)

	if !strings.HasPrefix(content, client.commandPrefix) {
		return
	}

	content = content[len(client.commandPrefix):]

	// split the command name with the arguments
	split := strings.SplitN(content, " ", -1)
	cmdName := strings.ToLower(split[0])
	cmdArgs := split[1:]

	for _, cmd := range commands.CommandList {
		if cmdName == cmd.Name {
			log.Println("Runnning commmand", cmdName, "with args", cmdArgs)

			client.wac.SendMessage(context.Background(), e.Info.Chat, &waProto.Message{
				ReactionMessage: &waProto.ReactionMessage{
					Key: &waProto.MessageKey{
						RemoteJid: proto.String(e.Info.Chat.String()),
						FromMe:    proto.Bool(true),
						Id:        &e.Info.ID,
					},
					Text:              proto.String("⏳"),
					SenderTimestampMs: proto.Int64(time.Now().UnixMilli()),
				},
			})

			err := cmd.Execute(client.wac, e, ctx, cmdArgs)

			var reaction string

			if err == nil {
				reaction = "👍"
			} else {
				log.Println(err)
				reaction = "❌"
			}

			client.wac.SendMessage(context.Background(), e.Info.Chat, &waProto.Message{
				ReactionMessage: &waProto.ReactionMessage{
					Key: &waProto.MessageKey{
						RemoteJid: proto.String(e.Info.Chat.String()),
						FromMe:    proto.Bool(true),
						Id:        &e.Info.ID,
					},
					Text:              &reaction,
					SenderTimestampMs: proto.Int64(time.Now().UnixMilli()),
				},
			})

			log.Println("Done.")
			break
		}
	}
}

func (client *PrinceClient) sendRepeatedMessages() {
	msgs := db.GetRepeatedMessageToday()

	for _, msg := range msgs {
		client.sendMessage(msg.JID, msg.Message)

		// Update next date
		switch msg.Repeat {
		case "y":
			msg.NextDate = msg.NextDate.AddDate(1, 0, 0)
		case "m":
			msg.NextDate = msg.NextDate.AddDate(0, 1, 0)
		case "w":
			msg.NextDate = msg.NextDate.AddDate(0, 0, 7)
		case "d":
			msg.NextDate = msg.NextDate.AddDate(0, 0, 1)
		}

		db.UpdateNextDate(msg.ID, msg.NextDate)
	}

	log.Println("Sent", len(msgs), "repeated messages")
}

func (client *PrinceClient) sendMessage(jidCode string, messageContent string) {
	jid, err := types.ParseJID(jidCode)
	if err != nil {
		log.Println("Failed to send message:", err)
		return
	}

	_, err = client.wac.SendMessage(context.Background(), jid, &waProto.Message{
		Conversation: proto.String(messageContent),
	})
	if err != nil {
		log.Println("Failed to send message:", err)
		return
	}
}
