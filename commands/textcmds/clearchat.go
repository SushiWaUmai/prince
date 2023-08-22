package textcmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func ClearChatCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	utils.ClearChat(chat)
	return nil, nil
}

func init() {
	utils.CreateCommand("clearchat", "ADMIN", "Clears the ChatGPT chat history", ClearChatCommand)
}
