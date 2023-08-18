package textcmds

import (
	"context"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func EchoCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	var msg *waProto.Message
	var text string
	if len(args) > 0 {
		text = strings.Join(args, " ")
		text = strings.TrimSpace(text)
		msg = utils.CreateTextMessage(text)
	} else {
		msg = pipe
	}

	_, err := client.SendMessage(context.Background(), chat, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func init() {
	utils.CreateCommand("echo", "USER", EchoCommand)
}
