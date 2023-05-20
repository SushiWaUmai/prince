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

	content, ctx := commands.GetTextContext(e.Message)

	if !strings.HasPrefix(content, client.commandPrefix) {
		return
	}

	content = content[len(client.commandPrefix):]
  content = strings.TrimSpace(content)

	commandsSplit := strings.Split(content, "|")
	commandInput := make([]commands.CommandInput, len(commandsSplit))
	for i, c := range commandsSplit {
		// split the command name with the arguments
		c = strings.TrimSpace(c)
		split := strings.Split(c, " ")
		commandInput[i] = commands.CommandInput{
			Name: strings.ToLower(split[0]),
			Args: split[1:],
		}
	}

	// Validate all commands
	for _, c := range commandInput {
		_, ok := commands.CommandMap[c.Name]
		if !ok {
			log.Println("Command not found: ", c.Name)
			return
		}
	}

	var pipe *waProto.Message = nil
	if ctx != nil {
		pipe = ctx.QuotedMessage
	}

	reaction := "⏳"

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

	var err error
	for _, c := range commandInput {
		log.Println("Runnning commmand", c.Name, "with args", c.Args)
		pipe, err = commands.CommandMap[c.Name].Execute(client.wac, e, ctx, pipe, c.Args)

		if err == nil {
			reaction = "👍"
		} else {
			log.Println(err)
			reaction = "❌"
			break
		}
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

	if pipe != nil {
		client.wac.SendMessage(context.Background(), e.Info.Chat, pipe)
	}

	log.Println("Done.")
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
