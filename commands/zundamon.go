package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SushiWaUmai/prince/env"
	"github.com/SushiWaUmai/prince/utils"
	"github.com/aethiopicuschan/voicevox"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func init() {
	voicevox := voicevox.NewClient("http", fmt.Sprintf("%s:50021", env.VOICEVOX_ENDPOINT))
	zundamonIdx := 1

	createCommand("zundamon", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		var text string
		if len(args) > 0 {
			text = strings.Join(args, " ")
		} else {
			text, _ = GetTextContext(pipe)
		}
		text = strings.TrimSpace(text)

		if text == "" {
			response := &waProto.Message{
				Conversation: proto.String("Please specify a text to speak"),
			}
			return response, errors.New("No Text specified")
		}

		query, err := voicevox.CreateQuery(zundamonIdx, text)

		if err != nil {
			return nil, err
		}

		wav, err := voicevox.CreateVoice(1, true, query)
		if err != nil {
			return nil, err
		}

		audioBytes, err := utils.ToOgg(wav)

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

		if err != nil {
			return nil, err
		}

		response := &waProto.Message{
			AudioMessage: audioMsg,
		}

		return response, nil
	})
}
