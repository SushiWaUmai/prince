package client

import (
	"log"
	"strings"
	"time"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func (client *PrinceClient) handleMessage(e *events.Message) {
	if !e.Info.IsFromMe {
		return
	}

	client.handleCommand(e.Message, e.Info.ID, e.Info.Chat)
}

func (client *PrinceClient) handleCommand(message *waProto.Message, msgId types.MessageID, chat types.JID) {
	content, ctx := utils.GetTextContext(message)

	if !strings.HasPrefix(content, client.commandPrefix) {
		return
	}

	content = content[len(client.commandPrefix):]
	content = strings.TrimSpace(content)

	commandsSplit := strings.Split(content, "|")
	commandInput := make([]utils.CommandInput, len(commandsSplit))
	for i, c := range commandsSplit {
		// split the command name with the arguments
		c = strings.TrimSpace(c)
		split := strings.Split(c, " ")
		commandInput[i] = utils.CommandInput{
			Name: strings.ToLower(split[0]),
			Args: split[1:],
		}
	}

	// Validate all commands
	for _, c := range commandInput {
		_, ok := utils.CommandMap[c.Name]
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

	client.SendMessage(chat, &waProto.Message{
		ReactionMessage: &waProto.ReactionMessage{
			Key: &waProto.MessageKey{
				RemoteJid: proto.String(chat.String()),
				FromMe:    proto.Bool(true),
				Id:        &msgId,
			},
			Text:              &reaction,
			SenderTimestampMs: proto.Int64(time.Now().UnixMilli()),
		},
	})

	var err error
	for _, c := range commandInput {
		log.Println("Runnning commmand", c.Name, "with args", c.Args)
		pipe, err = utils.CommandMap[c.Name].Execute(client.wac, chat, ctx, pipe, c.Args)

		if err == nil {
			reaction = "👍"
		} else {
			log.Println(err)
			reaction = "❌"
			break
		}
	}

	client.SendMessage(chat, &waProto.Message{
		ReactionMessage: &waProto.ReactionMessage{
			Key: &waProto.MessageKey{
				RemoteJid: proto.String(chat.String()),
				FromMe:    proto.Bool(true),
				Id:        &msgId,
			},
			Text:              &reaction,
			SenderTimestampMs: proto.Int64(time.Now().UnixMilli()),
		},
	})

	if pipe != nil {
		client.SendCommandMessage(chat, pipe)
	}

	log.Println("Done.")
}

func (client *PrinceClient) sendRepeatedMessages() {
	msgs := db.GetRepeatedMessageToday()

	for _, msg := range msgs {
		jid, err := types.ParseJID(msg.JID)
		if err != nil {
			log.Println("Failed to send message:", err)
			continue
		}
		_, err = client.SendMessage(jid, &waProto.Message{
			Conversation: proto.String(msg.Message),
		})
		if err != nil {
			log.Println("Failed to send message:", err)
			continue
		}

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
