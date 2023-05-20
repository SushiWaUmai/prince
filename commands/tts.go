package commands

import (
	"context"
	"errors"

	"os"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func init() {
	createCommand("tts", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		text, _ := GetTextContext(pipe)

		if text == "" {
			if len(args) <= 0 {
				response := &waProto.Message{
					Conversation: proto.String("Please specify a text to speak"),
				}
				return response, errors.New("No Text specified")
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
		audioBytes, err = utils.ToOgg(audioBytes)

		uploadResp, err := client.Upload(context.Background(), audioBytes, whatsmeow.MediaAudio)
		if err != nil {
			return nil, err
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

		response := &waProto.Message{
			AudioMessage: audioMsg,
		}
		return response, nil
	})
}
