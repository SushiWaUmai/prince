package mediacmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func MemeCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	subreddit := ""
	if len(args) > 0 {
		subreddit = args[0]
	}

	resp, err := utils.GetMeme(subreddit)
	if err != nil {
		return nil, err
	}

	buffer, err := utils.GetMemeImg(resp)
	if err != nil {
		return nil, err
	}

	response, err := utils.CreateImgMessage(client, buffer)
	if err != nil {
		return nil, err
	}

	response.ImageMessage.Caption = &resp.Title
	return response, nil
}

func init() {
	utils.CreateCommand("meme", "USER", "Retrieves and sends a meme from a specified subreddit", MemeCommand)
}
