package textcmds

import (
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func SystemChatCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	var text string
	if len(args) > 0 {
		text = strings.Join(args, " ")
	} else {
		text, _ = utils.GetTextContext(pipe)
	}
	text = strings.TrimSpace(text)

	utils.SetSystemChat(chat, text)

	return nil, nil
}

func init() {
	utils.CreateCommand("systemchat", "ADMIN", "Updates the system chat setting for ChatGPT", SystemChatCommand)
}
