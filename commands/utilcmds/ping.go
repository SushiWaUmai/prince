package utilcmds

import (
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func init() {
	utils.CreateCommand("ping", func(client *whatsmeow.Client, chat types.JID, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		reply := "pong!"
		if len(args) > 0 {
			reply = strings.Join(args, " ")
		}

		response := &waProto.Message{
			Conversation: &reply,
		}
		return response, nil
	})
}
