package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type animeResponse struct {
	URL string `json:"url"`
}

var animeCategories = []string{
	"waifu",
	"neko",
	"shinobu",
	"megumin",
	"bully",
	"cuddle",
	"cry",
	"hug",
	"awoo",
	"kiss",
	"lick",
	"pat",
	"smug",
	"bonk",
	"yeet",
	"blush",
	"smile",
	"wave",
	"highfive",
	"handhold",
	"nom",
	"bite",
	"glomp",
	"slap",
	"kill",
	"kick",
	"happy",
	"wink",
	"poke",
	"dance",
	"cringe",
}

var nsfwAnimeCategories = []string{
	"waifu",
	"neko",
	"trap",
	"blowjob",
}

func init() {
	createCommand("waifu", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) {
		category := "waifu"
		imgType := "sfw"

		// Check for arguments
		if len(args) > 0 {
			cLower := strings.ToLower(args[0])
			if contains(animeCategories, cLower) {
				category = cLower
			} else if contains(nsfwAnimeCategories, cLower) {
				category = cLower
				imgType = "nsfw"
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

		resp, err := http.Get(fmt.Sprintf("https://api.waifu.pics/%s/%s", imgType, category))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		var data animeResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			fmt.Println(err)
			return
		}

		resp, err = http.Get(data.URL)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		mimeType := resp.Header.Get("Content-Type")

		buffer, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		uploadResp, err := client.Upload(context.Background(), buffer, whatsmeow.MediaImage)
		if err != nil {
			fmt.Println(err)
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
