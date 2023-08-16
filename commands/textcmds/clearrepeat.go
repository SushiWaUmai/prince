package textcmds

import (
	"fmt"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
)

func ClearRepeatCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	// Delete the message
	affected := db.ClearRepeatedMessage(chat.String(), user)

	// Reply
	response := &waProto.Message{
		Conversation: proto.String(fmt.Sprintf("Deleted %d", affected)),
	}

	return response, nil
}

func init() {
	utils.CreateCommand("clearrepeat", "ADMIN", ClearRepeatCommand)
}
