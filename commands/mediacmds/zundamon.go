package mediacmds

import (
	"context"
	"errors"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

var voicevox = utils.CreateVoiceVoxClient()
var zundamonIdx = 1

func ZundamonCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
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

	response := &waProto.Message{
		AudioMessage: audioMsg,
	}

	return response, nil
}

func init() {
	utils.CreateCommand("zundamon", "USER", ZundamonCommand)
}
