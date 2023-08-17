package mediacmds

import (
	"errors"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func ToStickerCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	if pipe == nil || pipe.ImageMessage == nil {
		return utils.CreateTextMessage("Please reply to an image message"), errors.New("No ImageMessage quoted")
	}
	imgMsg := pipe.ImageMessage

	buffer, err := client.Download(imgMsg)
	if err != nil {
		return nil, err
	}
	response, err := utils.CreateStickerMessage(client, buffer)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func init() {
	utils.CreateCommand("tosticker", "USER", ToStickerCommand)
}
