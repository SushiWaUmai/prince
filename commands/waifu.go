package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
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
	utils.CreateCommand("waifu", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
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

				response := &waProto.Message{
					Conversation: &msg,
				}

				return response, nil
			}
		}

		resp, err := http.Get(fmt.Sprintf("https://api.waifu.im/search/?included_tags=%s", category))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var data animeResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}

		resp, err = http.Get(data.Images[0].Url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		mimeType := resp.Header.Get("Content-Type")

		buffer, err := ioutil.ReadAll(resp.Body)
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
