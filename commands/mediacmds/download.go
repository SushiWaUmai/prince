package mediacmds

import (
	"errors"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
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

		return utils.GetMedia(client, fetchUrl)
	})
}
