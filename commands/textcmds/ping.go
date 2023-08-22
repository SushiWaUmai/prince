package textcmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func PingCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	return utils.CreateTextMessage("pong!"), nil
}

func init() {
	utils.CreateCommand("ping", "USER", "Responds with the message \"pong!\" in the chat ", PingCommand)
}
