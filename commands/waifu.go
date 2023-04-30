package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type animeImage struct {
	Url string `json:"url"`
}

type animeResponse struct {
	Images []animeImage `json:"images"`
}

var animeCategories = []string{
	"maid",
	"waifu",
	"marin-kitagawa",
	"mori-calliope",
	"raiden-shogun",
	"oppai",
	"selfies",
	"uniform",
}

var nsfwAnimeCategories = []string{
	"ass",
	"hentai",
	"milf",
	"oral",
	"paizuri",
	"ecchi",
	"ero",
}

func init() {
	createCommand("waifu", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) {
		category := "waifu"

		// Check for arguments
		if len(args) > 0 {
			cLower := strings.ToLower(args[0])
			if contains(animeCategories, cLower) {
				category = cLower
			} else if contains(nsfwAnimeCategories, cLower) {
				category = cLower
			} else if cLower == "categories" {
				tLower := "sfw"

				if len(args) > 1 {
					tLower = strings.ToLower(args[1])
				}

				msg := "Categories:\n"
				if tLower == "nsfw" {
					msg += strings.Join(nsfwAnimeCategories, ", ")
				} else {
					msg += strings.Join(animeCategories, ", ")
				}

				client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
					Conversation: proto.String(msg),
				})
				return
			}
		}

		resp, err := http.Get(fmt.Sprintf("https://api.waifu.im/search/?included_tags=%s", category))
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		var data animeResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			log.Println(err)
			return
		}

		resp, err = http.Get(data.Images[0].Url)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		mimeType := resp.Header.Get("Content-Type")
		log.Println(mimeType)

		buffer, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		uploadResp, err := client.Upload(context.Background(), buffer, whatsmeow.MediaImage)
		if err != nil {
			log.Println(err)
			return
		}

		imgMsg := &waProto.ImageMessage{
			Mimetype:      proto.String(mimeType),
			Url:           &uploadResp.URL,
			DirectPath:    &uploadResp.DirectPath,
			MediaKey:      uploadResp.MediaKey,
			FileEncSha256: uploadResp.FileEncSHA256,
			FileSha256:    uploadResp.FileSHA256,
			FileLength:    &uploadResp.FileLength,
		}

		client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
			ImageMessage: imgMsg,
		})
	})
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
