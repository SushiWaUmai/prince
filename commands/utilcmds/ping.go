package utilcmds

import (
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func PingCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	reply := "pong!"
	if len(args) > 0 {
		reply = strings.Join(args, " ")
	}

	return utils.CreateTextMessage(reply), nil
}

func init() {
	utils.CreateCommand("ping", "USER", PingCommand)
}
