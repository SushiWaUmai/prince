package utilcmds

import (
	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func OnMessageClearCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	db.ClearMessageEvents(chat.String())
	return nil, nil
}

func init() {
	utils.CreateCommand("onmessageclear", "OP", "Clears all OnMessage events in this chat", OnMessageClearCommand)
}
