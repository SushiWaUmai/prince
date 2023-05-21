package mediacmds

import (
	"bytes"
	"context"
	"errors"
	"image"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"github.com/wader/goutubedl"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	"mvdan.cc/xurls/v2"
)

func init() {
	rxStrict := xurls.Strict()

	utils.CreateCommand("download", "USER", func(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		text, _ := utils.GetTextContext(pipe)
		text += " "

		text += strings.Join(args, " ")

		fetchUrl := rxStrict.FindString(text)

		if fetchUrl == "" {
			response := &waProto.Message{
				Conversation: proto.String("Please specify a url"),
			}
			return response, errors.New("No fetch url provoided")
		}

		resp, err := http.Get(fetchUrl)
		if err != nil {
			response := &waProto.Message{
				Conversation: proto.String("Failed fetch url"),
			}
			return response, errors.New("Failed to fetch url")
		}

		mimeType := resp.Header.Get("Content-Type")

		buffer, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if strings.Contains(mimeType, "image") {
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

			response := &waProto.Message{
				ImageMessage: imgMsg,
			}
			return response, nil
		} else if strings.Contains(mimeType, "audio") {
			uploadResp, err := client.Upload(context.Background(), buffer, whatsmeow.MediaAudio)
			if err != nil {
				return nil, err
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

			response := &waProto.Message{
				AudioMessage: audioMsg,
			}
			return response, nil
		} else {
			// yt-dlp
			goutubedl.Path = "yt-dlp"
			result, err := goutubedl.New(context.Background(), fetchUrl, goutubedl.Options{})
			if err != nil {
				return nil, err
			}
			downloadResult, err := result.Download(context.Background(), "best")
			if err != nil {
				return nil, err
			}
			defer downloadResult.Close()

			buffer, err := ioutil.ReadAll(downloadResult)
			if err != nil {
				return nil, err
			}

			uploadResp, err := client.Upload(context.Background(), buffer, whatsmeow.MediaVideo)
			if err != nil {
				return nil, err
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

			response := &waProto.Message{
				VideoMessage: videoMsg,
			}
			return response, nil
		}
	})
}
