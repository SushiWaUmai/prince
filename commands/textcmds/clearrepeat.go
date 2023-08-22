package textcmds

import (
	"fmt"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
)

// TODO: Convert repeated messages to repeated commands
func ClearRepeatCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	// Delete the message
	affected, err := db.ClearRepeatedMessage(chat.String(), user)
	if err != nil {
		return nil, err
	}

	return utils.CreateTextMessage(fmt.Sprintf("Deleted %d", affected)), nil
}

func init() {
	utils.CreateCommand("clearrepeat", "ADMIN", "", ClearRepeatCommand)
}
