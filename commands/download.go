package commands

import (
	"bytes"
	"context"
	"errors"
	"image"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/wader/goutubedl"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
	"mvdan.cc/xurls/v2"
)

func init() {
	rxStrict := xurls.Strict()

	createCommand("download", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) error {
		var text string

		if ctx != nil && ctx.QuotedMessage != nil {
			if ctx.QuotedMessage.Conversation != nil {
				text = *ctx.QuotedMessage.Conversation + " "
			}

			if ctx.QuotedMessage.ExtendedTextMessage != nil && ctx.QuotedMessage.ExtendedTextMessage.Text != nil {
				text = *ctx.QuotedMessage.ExtendedTextMessage.Text + " "
			}
		}

		text += strings.Join(args, " ")

		fetchUrl := rxStrict.FindString(text)

		if fetchUrl == "" {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Please specify a url"),
			})
			return errors.New("No fetch url provoided")
		}

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
		} else {
			// yt-dlp
			goutubedl.Path = "yt-dlp"
			result, err := goutubedl.New(context.Background(), fetchUrl, goutubedl.Options{})
			if err != nil {
				return err
			}
			downloadResult, err := result.Download(context.Background(), "best")
			if err != nil {
				return err
			}
			defer downloadResult.Close()

			buffer, err := ioutil.ReadAll(downloadResult)
			if err != nil {
				return err
			}

			uploadResp, err := client.Upload(context.Background(), buffer, whatsmeow.MediaVideo)
			if err != nil {
				return err
			}

			videoMsg := &waProto.VideoMessage{
				Mimetype:      proto.String(http.DetectContentType(buffer)),
				Url:           &uploadResp.URL,
				DirectPath:    &uploadResp.DirectPath,
				MediaKey:      uploadResp.MediaKey,
				FileEncSha256: uploadResp.FileEncSHA256,
				FileSha256:    uploadResp.FileSHA256,
				FileLength:    &uploadResp.FileLength,
			}

			_, err = client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				VideoMessage: videoMsg,
			})
		}

		if err != nil {
			return err
		}

		return nil
	})
}
