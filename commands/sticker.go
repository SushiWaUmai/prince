package commands

import (
	"bytes"
	"context"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"log"

	"github.com/chai2010/webp"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func init() {
	createCommand("sticker", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) {
		if ctx == nil || ctx.QuotedMessage == nil || ctx.QuotedMessage.ImageMessage == nil {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Please reply to a image message"),
			})
			return
		}
		imgMsg := ctx.QuotedMessage.ImageMessage

		buffer, err := client.Download(imgMsg)
		if err != nil {
			log.Println(err)
			return
		}

		img, _, err := image.Decode(bytes.NewReader(buffer))
		if err != nil {
			log.Fatal(err)
		}
		g := img.Bounds()

		// Get height and width
		width := uint32(g.Dx())
		height := uint32(g.Dy())

		webpByte, err := webp.EncodeRGBA(img, *proto.Float32(1))
		if err != nil {
			log.Fatal(err)
			return
		}

		uploadResp, err := client.Upload(context.Background(), webpByte, whatsmeow.MediaImage)
		if err != nil {
			log.Println(err)
			return
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

		client.SendMessage(
			context.Background(),
			messageEvent.Info.Chat,
			&waProto.Message{
				StickerMessage: stickerMsg,
			},
		)
	})
}
