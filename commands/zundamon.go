package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SushiWaUmai/prince/env"
	"github.com/aethiopicuschan/voicevox"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func init() {
	voicevox := voicevox.NewClient("http", fmt.Sprintf("%s:50021", env.VOICEVOX_ENDPOINT))
	zundamonIdx := 1

	createCommand("zundamon", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) error {
		var text string
		pipeString, _ := GetTextContext(pipe)

		if pipeString != "" {
			text = pipeString
		} else {
			if len(args) <= 0 {
				client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
					Conversation: proto.String("Please specify a text to speak"),
				})
				return errors.New("No Text specified")
			}

			text = strings.Join(args, " ")
		}

		query, err := voicevox.CreateQuery(zundamonIdx, text)

		if err != nil {
			return err
		}

		wav, err := voicevox.CreateVoice(1, true, query)
		if err != nil {
			return err
		}

		audioBytes, err := toOgg(wav)

		uploadResp, err := client.Upload(context.Background(), audioBytes, whatsmeow.MediaAudio)
		if err != nil {
			return err
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

		if err != nil {
			return err
		}

		_, err = client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
			AudioMessage: audioMsg,
		})

		return nil
	})
}
