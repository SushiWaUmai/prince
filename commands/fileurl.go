package commands

import (
	"bytes"
	"context"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func init() {
	createCommand("fileurl", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) {
		if len(args) <= 0 {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Please specify a fetchUrl"),
			})
			return
		}

		fetchUrl := args[0]

		resp, err := http.Get(fetchUrl)
		if err != nil {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Failed fetch Url"),
			})
			return
		}

		mimeType := resp.Header.Get("Content-Type")

		buffer, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		if strings.Contains(mimeType, "image") {
			uploadResp, err := client.Upload(context.Background(), buffer, whatsmeow.MediaImage)
			if err != nil {
				log.Println(err)
				return
			}

			img, _, err := image.Decode(bytes.NewBuffer(buffer))
			if err != nil {
				log.Fatal(err)
			}
			g := img.Bounds()

			// Get height and width
			width := uint32(g.Dx())
			height := uint32(g.Dy())

			imgMsg := &waProto.ImageMessage{
				Mimetype:      &mimeType,
				Url:           &uploadResp.URL,
				DirectPath:    &uploadResp.DirectPath,
				MediaKey:      uploadResp.MediaKey,
				FileEncSha256: uploadResp.FileEncSHA256,
				FileSha256:    uploadResp.FileSHA256,
				FileLength:    &uploadResp.FileLength,
				Width:         &width,
				Height:        &height,
			}

			_, err = client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				ImageMessage: imgMsg,
			})

		} else if strings.Contains(mimeType, "audio") {
			uploadResp, err := client.Upload(context.Background(), buffer, whatsmeow.MediaAudio)
			if err != nil {
				log.Println(err)
				return
			}

			audioMsg := &waProto.AudioMessage{
				Mimetype:      &mimeType,
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

		}

		if err != nil {
			log.Println(err)
		}
	})
}
