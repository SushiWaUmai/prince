package mediacmds

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/SushiWaUmai/prince/env"
	"github.com/SushiWaUmai/prince/utils"
	"github.com/aethiopicuschan/voicevox"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func init() {
	// TODO: Move this to utils
	url, err := url.Parse(env.VOICEVOX_ENDPOINT)
	if err != nil {
		log.Println("Failed to register zundamon command:", err)
	}

	voicevox := voicevox.NewClient(url.Scheme, url.Host)
	zundamonIdx := 1

	utils.CreateCommand("zundamon", func(client *whatsmeow.Client, chat types.JID, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		var text string
		if len(args) > 0 {
			text = strings.Join(args, " ")
		} else {
			text, _ = utils.GetTextContext(pipe)
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
