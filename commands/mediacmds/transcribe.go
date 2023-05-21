package mediacmds

import (
	"errors"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func init() {
	utils.CreateCommand("transcribe", func(client *whatsmeow.Client, chat types.JID, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		// Check if there's a voice message quoted
		if pipe.AudioMessage == nil {
			response := &waProto.Message{
				Conversation: proto.String("Please reply to a voice message"),
			}
			return response, errors.New("No voice message quoted")
		}

		// Download the voice message
		audioData, err := client.Download(pipe.AudioMessage)
		if err != nil {
			response := &waProto.Message{
				Conversation: proto.String("Failed to download the voice message"),
			}
			return response, err
		}

		// Use a Golang library to transcribe the audio to text
		transcription, err := utils.TranscribeAudio(audioData)
		if err != nil {
			response := &waProto.Message{
				Conversation: proto.String("Failed to transcribe the audio file"),
			}
			return response, err
		}

		// Send the transcription back to the user
		response := &waProto.Message{
			Conversation: &transcription,
		}
		return response, nil
	})
}
