package mediacmds

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
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

		imgMsg, err := utils.CreateImgMessage(client, buffer)
		if err != nil {
			return nil, err
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
