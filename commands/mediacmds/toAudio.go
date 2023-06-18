package mediacmds

import (
	"context"
	"errors"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func init() {
	utils.CreateCommand("toaudio", "USER", func(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		if pipe == nil || pipe.VideoMessage == nil {
			response := &waProto.Message{
				Conversation: proto.String("Please reply to a video message"),
			}
			return response, errors.New("No VideoMessage quoted")
		}
		videoMessage := pipe.VideoMessage

		buffer, err := client.Download(videoMessage)
		if buffer == nil {
			return nil, err
		}

		audioBytes, err := utils.VideoToAudio(buffer)

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
