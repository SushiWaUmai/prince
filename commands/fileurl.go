package commands

import (
	"bytes"
	"context"
	"errors"
	"image"
	"io/ioutil"
	"net/http"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
	"mvdan.cc/xurls/v2"
)

var rxStrict = xurls.Strict()

func init() {
	createCommand("fileurl", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) error {
		var text string

		if ctx != nil && ctx.QuotedMessage != nil && ctx.QuotedMessage.Conversation != nil {
			text = *ctx.QuotedMessage.Conversation + " "
		}

		text += strings.Join(args, " ")
		urls := rxStrict.FindAllString(text, -1)

		if len(urls) <= 0 {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Please specify a url"),
			})
			return errors.New("No fetch url provoided")
		}
		fetchUrl := urls[0]

		resp, err := http.Get(fetchUrl)
		if err != nil {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Failed fetch url"),
			})
			return errors.New("Failed to fetch url")
		}

		mimeType := resp.Header.Get("Content-Type")

		buffer, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if strings.Contains(mimeType, "image") {
			uploadResp, err := client.Upload(context.Background(), buffer, whatsmeow.MediaImage)
			if err != nil {
				return err
			}

			img, _, err := image.Decode(bytes.NewBuffer(buffer))
			if err != nil {
				return err
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
				return err
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
			return err
		}

		return nil
	})
}
