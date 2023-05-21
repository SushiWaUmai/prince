package mediacmds

import (
	"bytes"
	"context"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

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
	utils.CreateCommand("waifu", "USER", func(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
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

		buffer, err := utils.GetWaifu(category)
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
			Mimetype:      proto.String(http.DetectContentType(buffer)),
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
