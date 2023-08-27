package client

import (
	"log"
	"strings"
	"time"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/lang"
	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (client *PrinceClient) handleCommand(message *waProto.Message, msgId types.MessageID, chat types.JID, user string) {
	content, ctx := utils.GetTextContext(message)

	if !strings.HasPrefix(content, string(client.commandPrefix)) {
		return
	}

	fromMe := client.wac.Store.ID.ToNonAD().User == user
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

	result, err := RunCommand(client.wac, content, ctx, chat, user)

	if err == nil {
		reaction = "👍"
	} else {
		log.Println(err)
		reaction = "❌"
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

	if result != nil {
		client.SendCommandMessage(chat, user, result)
	}

	log.Println("Done.")
}

func RunCommand(client *whatsmeow.Client, content string, ctx *waProto.ContextInfo, chat types.JID, user string) (*waProto.Message, error) {
	userPermission, err := db.GetUserPermission(user)
	if err != nil {
		return nil, err
	}

	perm := userPermission.Permission
	fromMe := client.Store.ID.ToNonAD().User == user
	if fromMe {
		perm = "OP"
	}

	if perm == "NONE" {
		return nil, utils.ErrNotEnoughPermission
	}

	commandInput, err := lang.Scan(content)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Validate all commands
	for _, c := range commandInput {
		cmd, _ := utils.CommandMap[c.Name]
		if !db.ComparePermission(perm, cmd.Permission) {
			return utils.CreateTextMessage("You do not have enough permission to run this command."), utils.ErrNotEnoughPermission
		}
	}

	var pipe *waProto.Message = nil
	if ctx != nil {
		pipe = ctx.QuotedMessage
	}

	for _, c := range commandInput {
		log.Println("Runnning commmand", c.Name, "with args", c.Args)
		pipe, err = utils.CommandMap[c.Name].Execute(client, chat, user, ctx, pipe, c.Args)

		if err != nil {
			return nil, err
		}
	}

	return pipe, nil
}
