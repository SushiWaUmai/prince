package commands

import (
	"bytes"
	"context"
	"io"
	"log"

	"os"
	"strings"

	"github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func init() {
	createCommand("tts", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) {
		var text string

		if ctx != nil && ctx.QuotedMessage != nil && ctx.QuotedMessage.Conversation != nil {
			text = *ctx.QuotedMessage.Conversation
		} else {
			if len(args) <= 0 {
				client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
					Conversation: proto.String("Please specify a text to speak"),
				})
				return
			}

			text = strings.Join(args, " ")
		}

		speech := htgotts.Speech{
			Folder:   ".",
			Language: voices.English,
		}

		speech.CreateSpeechFile(text, "speach")
		defer os.Remove("speach.mp3")

		// get the bytes
		audioBytes, err := os.ReadFile("speach.mp3")
		audioBytes, err = mp3ToOgg(audioBytes)

		uploadResp, err := client.Upload(context.Background(), audioBytes, whatsmeow.MediaAudio)
		if err != nil {
			log.Println(err)
			return
		}

		audioMsg := &waProto.AudioMessage{
			Mimetype:      proto.String("audio/ogg; codecs=opus"),
			Url:           &uploadResp.URL,
			DirectPath:    &uploadResp.DirectPath,
			MediaKey:      uploadResp.MediaKey,
			FileEncSha256: uploadResp.FileEncSHA256,
			FileSha256:    uploadResp.FileSHA256,
			FileLength:    &uploadResp.FileLength,
		}

		_, err = client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
			AudioMessage: audioMsg,
		})

		if err != nil {
			log.Println(err)
		}
	})
}

func mp3ToOgg(audioData []byte) ([]byte, error) {
	// ffmpeg -i $inFileName -acodec libmp3lame -y $outFileName
	tmpFile, err := os.CreateTemp("", "audio*.mp3")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, bytes.NewReader(audioData))
	if err != nil {
		return nil, err
	}

	tmpFileOut, err := os.CreateTemp("", "audio*.ogg")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFileOut.Name())

	err = ffmpeg.Input(tmpFile.Name()).Output(tmpFileOut.Name(), ffmpeg.KwArgs{
		"acodec": "libmp3lame",
		"c:a": "libopus",
	}).OverWriteOutput().Run()

	if err != nil {
		return nil, err
	}

	return os.ReadFile(tmpFileOut.Name())
}
