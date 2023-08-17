package mediacmds

import (
	"errors"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func ToAudioCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	if pipe == nil || pipe.VideoMessage == nil {
		return utils.CreateTextMessage("Please reply to a video message"), errors.New("No VideoMessage quoted")
	}
	videoMessage := pipe.VideoMessage

	buffer, err := client.Download(videoMessage)
	if buffer == nil {
		return nil, err
	}

	audioBytes, err := utils.VideoToAudio(buffer)
	if err != nil {
		return nil, err
	}

	response, err := utils.CreateAudioMessage(client, audioBytes)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func init() {
	utils.CreateCommand("toaudio", "USER", ToAudioCommand)
}
