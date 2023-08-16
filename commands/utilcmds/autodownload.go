package utilcmds

import (
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
)

func AutoDownloadCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	enabled := db.ToggleMessageEvent(chat.String(), "DOWNLOAD")
	var reply string

	if enabled {
		reply = "AutoDownload enabled"
	} else {
		reply = "AutoDownload disabled"
	}

	response := &waProto.Message{
		Conversation: &reply,
	}

	return response, nil
}

func init() {
	utils.CreateCommand("autodownload", "OP", AutoDownloadCommand)
}
