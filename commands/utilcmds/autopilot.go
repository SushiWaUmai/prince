package utilcmds

import (
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
)

func AutoPilotCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	enabled, err := db.ToggleMessageEvent(chat.String(), "CHAT")
	if err != nil {
		return nil, err
	}

	if enabled {
		return utils.CreateTextMessage("AutoPilot enabled"), nil
	} else {
		return utils.CreateTextMessage("AutoPilot disabled"), nil
	}
}

func init() {
	utils.CreateCommand("autopilot", "OP", AutoPilotCommand)
}
