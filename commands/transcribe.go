package commands

import (
	"context"
	"errors"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func init() {
	createCommand("transcribe", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) error {
		// Check if there's a voice message quoted
		if pipe.AudioMessage == nil {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Please reply to a voice message"),
			})
			return errors.New("No voice message quoted")
		}

		// Download the voice message
		audioData, err := client.Download(pipe.AudioMessage)
		if err != nil {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Failed to download the voice message"),
			})
			return err
		}

		// Use a Golang library to transcribe the audio to text
		transcription, err := utils.TranscribeAudio(audioData)
		if err != nil {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Failed to transcribe the audio file"),
			})
			return err
		}

		// Send the transcription back to the user
		_, err = client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
			Conversation: &transcription,
		})

		if err != nil {
			return err
		}

		return nil
	})
}
