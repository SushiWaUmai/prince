package utilcmds

import (
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
)

func init() {
	utils.CreateCommand("autopilot", "OP", func(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		enabled := db.ToggleMessageEvent(chat.String(), "CHAT")
		var reply string

		if enabled {
			reply = "AutoPilot enabled"
		} else {
			reply = "AutoPilot disabled"
		}

		response := &waProto.Message{
			Conversation: &reply,
		}

		return response, nil
	})
}

