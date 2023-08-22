package mediacmds

import (
	"errors"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func TranscribeCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	// Check if there's a voice message quoted
	if pipe.AudioMessage == nil {
		return utils.CreateTextMessage("Please reply to a voice message"), errors.New("No voice message quoted")
	}

	// Download the voice message
	audioData, err := client.Download(pipe.AudioMessage)
	if err != nil {
		return utils.CreateTextMessage("Failed to download the voice message"), err
	}

	// Use a Golang library to transcribe the audio to text
	transcription, err := utils.TranscribeAudio(audioData)
	if err != nil {
		return utils.CreateTextMessage("Failed to transcribe the audio file"), err
	}

	// Send the transcription back to the user
	return utils.CreateTextMessage(transcription), nil
}

func init() {
	utils.CreateCommand("transcribe", "ADMIN", "Converts a voice message into text",TranscribeCommand)
}
