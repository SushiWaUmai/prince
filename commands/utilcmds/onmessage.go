package utilcmds

import (
	"strings"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func OnMessageCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	var content string
	if len(args) > 0 {
		content = strings.Join(args, " ")
	} else {
		content, _ = utils.GetTextContext(pipe)
	}
	content = strings.TrimSpace(content)

	if content == "" {
		// List all MessageEvents
		events := db.GetMessageEvents(chat.String())

		if len(events) == 0 {
			return utils.CreateTextMessage("No MessageEvents registered"), nil
		}

		var msg []string
		msg = append(msg, "MessageEvents:")

		for _, event := range events {
			msg = append(msg, event.Content)
		}

		return utils.CreateTextMessage(strings.Join(msg, "\n")), nil
	}

	_, err := db.CreateMessageEvent(chat.String(), content)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func init() {
	utils.CreateCommand("onmessage", "OP", "Creates an OnMessage event with the given command in this chat", OnMessageCommand)
}
