package client

import (
	"log"
	"strings"
	"time"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (client *PrinceClient) handleCommand(message *waProto.Message, msgId types.MessageID, chat types.JID, user string) {
	content, ctx := utils.GetTextContext(message)

	if !strings.HasPrefix(content, client.commandPrefix) {
		return
	}

	perm := db.GetUserPermission(user).Permission
	fromMe := client.wac.Store.ID.ToNonAD().User == user
	if fromMe {
		perm = "OP"
	}

	if perm == "NONE" {
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
		cmd, ok := utils.CommandMap[c.Name]
		if !ok {
			log.Println("Command not found: ", c.Name)
			return
		}
		if !db.ComparePermission(perm, cmd.Permission) {
			log.Println("Not enough permission")
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
				FromMe:    &fromMe,
				Id:        &msgId,
			},
			Text:              &reaction,
			SenderTimestampMs: proto.Int64(time.Now().UnixMilli()),
		},
	})

	var err error
	for _, c := range commandInput {
		log.Println("Runnning commmand", c.Name, "with args", c.Args)
		pipe, err = utils.CommandMap[c.Name].Execute(client.wac, chat, user, ctx, pipe, c.Args)

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
				FromMe:    &fromMe,
				Id:        &msgId,
			},
			Text:              &reaction,
			SenderTimestampMs: proto.Int64(time.Now().UnixMilli()),
		},
	})

	if pipe != nil {
		client.SendCommandMessage(chat, user, pipe)
	}

	log.Println("Done.")
}
