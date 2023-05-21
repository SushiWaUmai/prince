package mediacmds

import (
	"bytes"
	"context"
	"image"
	"net/http"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func init() {
	utils.CreateCommand("sd", func(client *whatsmeow.Client, chat types.JID, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		var prompt string
		pipeString, _ := utils.GetTextContext(pipe)
		if pipeString != "" {
			prompt = pipeString
		}

		if len(args) > 0 {
			if prompt != "" {
				prompt += ", "
			}

			prompt += strings.Join(args, " ")
		}

		buffer, err := utils.Txt2Img(prompt)
		if err != nil {
			return nil, err
		}

		uploadResp, err := client.Upload(context.Background(), buffer, whatsmeow.MediaImage)
		if err != nil {
			return nil, err
		}

		img, _, err := image.Decode(bytes.NewBuffer(buffer))
		if err != nil {
			return nil, err
		}
		g := img.Bounds()

		// Get height and width
		width := uint32(g.Dx())
		height := uint32(g.Dy())

		imgMsg := &waProto.ImageMessage{
			Mimetype:      proto.String(http.DetectContentType(buffer)),
			Url:           &uploadResp.URL,
			DirectPath:    &uploadResp.DirectPath,
			MediaKey:      uploadResp.MediaKey,
			FileEncSha256: uploadResp.FileEncSHA256,
			FileSha256:    uploadResp.FileSHA256,
			FileLength:    &uploadResp.FileLength,
			Width:         &width,
			Height:        &height,
		}

		response := &waProto.Message{
			ImageMessage: imgMsg,
		}
		return response, nil
	})
}
