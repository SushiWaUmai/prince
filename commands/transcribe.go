package commands

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func init() {
	createCommand("transcribe", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) {
		// Check if there's a voice message quoted
		if ctx == nil || ctx.QuotedMessage == nil || ctx.QuotedMessage.AudioMessage == nil {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Please reply to a voice message"),
			})
			return
		}

		// Download the voice message
		audioData, err := client.Download(ctx.QuotedMessage.AudioMessage)
		if err != nil {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Failed to download the voice message"),
			})
			return
		}

		// Use a Golang library to transcribe the audio to text
		transcription, err := TranscribeAudio(audioData)
		if err != nil {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Failed to transcribe the audio file"),
			})
			log.Println(err)
			return
		}

		// Send the transcription back to the user
		client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
			Conversation: &transcription,
		})
	})
}

func TranscribeAudio(audioData []byte) (string, error) {
	audioData, err := oggToMp3(audioData)
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", "audio*.mp3")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, bytes.NewReader(audioData))
	if err != nil {
		return "", err
	}

	req, err := OpenAIClient.CreateTranscription(
		context.Background(),
		openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: tmpFile.Name(),
		},
	)
	if err != nil {
		return "", err
	}

	return req.Text, nil
}

func oggToMp3(audioData []byte) ([]byte, error) {
	// ffmpeg -i $inFileName -acodec libmp3lame -y $outFileName
	tmpFile, err := os.CreateTemp("", "audio*.ogg")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, bytes.NewReader(audioData))
	if err != nil {
		return nil, err
	}

	tmpFileOut, err := os.CreateTemp("", "audio*.mp3")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFileOut.Name())

	err = ffmpeg.Input(tmpFile.Name()).Output(tmpFileOut.Name(), ffmpeg.KwArgs{
		"acodec": "libmp3lame",
	}).OverWriteOutput().Run()

	if err != nil {
		return nil, err
	}

	return os.ReadFile(tmpFileOut.Name())
}
