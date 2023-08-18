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
	if len(args) > 0 {
		text := strings.Join(args, " ")
		text = strings.TrimSpace(text)
		return utils.CreateTextMessage(text), nil
	}

	_, err := client.SendMessage(context.Background(), chat, pipe)
	if err != nil {
		return nil, err
	}

	return pipe, nil
}

func init() {
	utils.CreateCommand("echo", "USER", EchoCommand)
}
