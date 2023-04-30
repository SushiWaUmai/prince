package commands

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/srwiley/oksvg"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func todoinit() {
	createCommand("github", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) {
		// get username
		if len(args) == 0 {
			client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
				Conversation: proto.String("Please provide a username"),
			})
			return
		}

		username := args[0]

		resp, err := http.Get(fmt.Sprintf("https://github-profile-summary-cards.vercel.app/api/cards/profile-details?username=%s", username))
		if err != nil {
			log.Println(err)
			return
		}

		svgBytes, err := ioutil.ReadAll(resp.Body)
		svgStream := bytes.NewBuffer(svgBytes)

		uploadResp, err := client.Upload(context.Background(), svgStream.Bytes(), whatsmeow.MediaImage)
		if err != nil {
			log.Println(err)
			return
		}

		icon, _ := oksvg.ReadIconStream(svgStream)
		width := uint32(icon.ViewBox.W)
		height := uint32(icon.ViewBox.H)

		mimeType := resp.Header.Get("Content-Type")

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

		if err != nil {
			log.Println(err)
		}
	})
}
