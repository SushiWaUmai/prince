package client

import (
	"log"
	"strings"
	"time"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/lang"
	"github.com/SushiWaUmai/prince/utils"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (client *PrinceClient) handleCommand(message *waProto.Message, msgId types.MessageID, chat types.JID, user string, silent bool) {
	content, ctx := utils.GetTextContext(message)

	if !strings.HasPrefix(content, string(client.commandPrefix)) {
		return
	}

	userPermission, err := db.GetUserPermission(user)
	if err != nil {
		return
	}

	perm := userPermission.Permission
	fromMe := client.wac.Store.ID.ToNonAD().User == user
	if fromMe {
		perm = "OP"
	}

	if perm == "NONE" {
		return
	}

	commandInput, err := lang.Scan(content)
	if err != nil {
		log.Println(err)
		return
	}

	// Validate all commands
	for _, c := range commandInput {
		cmd, _ := utils.CommandMap[c.Name]
		if !db.ComparePermission(perm, cmd.Permission) {
			if !silent {
				client.SendMessage(chat, utils.CreateTextMessage("You do not have enough permission to run this command."))
			}
			log.Println("Not enough permission")
			return
		}
	}

	var pipe *waProto.Message = nil
	if ctx != nil {
		pipe = ctx.QuotedMessage
	}

	reaction := "⏳"

	if !silent {
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
	}

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

	if !silent {
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

	}
	if pipe != nil && (!silent || err == nil){
		client.SendCommandMessage(chat, user, pipe)
	}

	log.Println("Done.")
}
