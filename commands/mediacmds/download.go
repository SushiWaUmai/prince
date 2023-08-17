package mediacmds

import (
	"errors"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"mvdan.cc/xurls/v2"
)

var rxStrict = xurls.Strict()

func DownloadCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	text, _ := utils.GetTextContext(pipe)
	text += " "

	text += strings.Join(args, " ")

	fetchUrl := rxStrict.FindString(text)

	if fetchUrl == "" {
		return utils.CreateTextMessage("Please specify a url"), errors.New("No fetch url provided")
	}

	return utils.GetMedia(client, fetchUrl)

}

func init() {
	utils.CreateCommand("download", "USER", DownloadCommand)
}
