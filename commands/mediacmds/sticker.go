package mediacmds

import (
	"bytes"
	"context"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/SushiWaUmai/prince/utils"
	"github.com/chai2010/webp"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func init() {
	utils.CreateCommand("sticker", func(client *whatsmeow.Client, chat types.JID, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		if pipe == nil || pipe.ImageMessage == nil {
			response := &waProto.Message{
				Conversation: proto.String("Please reply to a image message"),
			}
			return response, errors.New("No ImageMessage quoted")
		}
		imgMsg := pipe.ImageMessage

		buffer, err := client.Download(imgMsg)
		if err != nil {
			return nil, err
		}

		img, _, err := image.Decode(bytes.NewReader(buffer))
		if err != nil {
			return nil, err
		}
		g := img.Bounds()

		// Get height and width
		width := uint32(g.Dx())
		height := uint32(g.Dy())

		webpByte, err := webp.EncodeRGBA(img, *proto.Float32(1))
		if err != nil {
			return nil, err
		}

		uploadResp, err := client.Upload(context.Background(), webpByte, whatsmeow.MediaImage)
		if err != nil {
			return nil, err
		}

		stickerMsg := &waProto.StickerMessage{
			Mimetype:      proto.String(http.DetectContentType(webpByte)),
			Url:           &uploadResp.URL,
			DirectPath:    &uploadResp.DirectPath,
			MediaKey:      uploadResp.MediaKey,
			FileEncSha256: uploadResp.FileEncSHA256,
			FileSha256:    uploadResp.FileSHA256,
			FileLength:    &uploadResp.FileLength,
			PngThumbnail:  webpByte,
			Width:         &width,
			Height:        &height,
		}

		response := &waProto.Message{
			StickerMessage: stickerMsg,
		}
		return response, nil
	})
}
