package mediacmds

import (
	"errors"
	"log"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"mvdan.cc/xurls/v2"
)

var rxStrict = xurls.Strict()

func DownloadCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	text := strings.Join(args, " ")
	content, _ := utils.GetTextContext(pipe)
	text += " "
	text += content

	fetchUrl := rxStrict.FindString(text)
	log.Println(text)
	log.Println(fetchUrl)

	if fetchUrl == "" {
		return utils.CreateTextMessage("Please specify a url"), errors.New("No fetch url provided")
	}

	return utils.GetMedia(client, fetchUrl)
}

func init() {
	utils.CreateCommand("download", "USER", "Downloads and sends media from a specified URL in a chat message.", DownloadCommand)
}
